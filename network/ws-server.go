package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// wsTCPKeepAliveListener : 重载net/http tcpKeepAliveListener
type wsTCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept : 重载net/http wsTCPKeepAliveListener.Accept
func (ln wsTCPKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type wsServer struct {
	sServer
	waccept websocket.Upgrader
	httpSrv *http.Server
	httpMtx *http.ServeMux
	httpErr error
}

func (wss *wsServer) listen(operator *actor.PID, addr string) error {
	wss.waccept = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wss.maker = wss.wsMakeConn
	wss.httpMtx = http.NewServeMux()
	wss.httpSrv = &http.Server{Addr: addr, Handler: wss.httpMtx}
	wss.httpMtx.HandleFunc("/ws", wss.wsAccept)

	ln, err := wss.httpListen(addr)
	if err != nil {
		return err
	}

	wss.netWait.Add(2)
	go wss.serve(ln)
	go wss.keeploop()

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (wss *wsServer) serve(ln net.Listener) {
	defer wss.netWait.Done()
	for {
		wss.httpErr = wss.httpSrv.Serve(wsTCPKeepAliveListener{ln.(*net.TCPListener)})
		wss.isShutdown = true
		break
	}
}

func (wss *wsServer) wsAccept(w http.ResponseWriter, r *http.Request) {

	if wss.stat != Connected {
		return
	}

	s, err := wss.waccept.Upgrade(w, r, nil)
	if err != nil {
		//错误日志
		logger.Fatal(wss.operator.ID, "web socket accept fail:%v", err)
		return
	}

	err = wss.accept(s, s.RemoteAddr().Network(), s.RemoteAddr().String())
	if err != nil {
		logger.Fatal(wss.operator.ID, "web socket accept fail:%v", err)
	}
}

func (wss *wsServer) wsMakeConn(handle int32, s interface{}, operator *actor.PID, so *slot, now uint64, stat int32) ISocket {
	conn := &wsConn{}
	conn.h = handle
	conn.s = s
	conn.o = operator
	conn.so = so
	conn.stat = stat
	conn.rv = wsConnRecv
	conn.wr = wsConnWrite
	conn.cls = wsConnClose
	conn.out = make(chan *NetChunk, wss.outChanMax)
	conn.i.ReadLastTime = now
	conn.i.WriteLastTime = now

	return conn
}

func (wss *wsServer) httpListen(addr string) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}

	return net.Listen("tcp", addr)
}

func (wss *wsServer) getProto() string {
	return ProtoWeb
}
func (wss *wsServer) getType() int {
	return CListen
}

func (wss *wsServer) close(lck *util.ReSpinLock) {
	wss.httpSrv.Close()
}
