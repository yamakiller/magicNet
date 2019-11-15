package network

import (
	"net"
	"net/http"
	"time"

	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"

	"github.com/gorilla/websocket"
)

// wsTCPKeepAliveListener : 重载net/http tcpKeepAliveListener
type wsTCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept : 重载net/http wsTCPKeepAliveListener.Accept
func (slf wsTCPKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := slf.AcceptTCP()
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

func (slf *wsServer) listen(operator *actor.PID, addr string) error {
	slf.waccept = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	slf.maker = slf.wsMakeConn
	slf.httpMtx = http.NewServeMux()
	slf.httpSrv = &http.Server{Addr: addr, Handler: slf.httpMtx}
	slf.httpMtx.HandleFunc("/ws", slf.wsAccept)

	ln, err := slf.httpListen(addr)
	if err != nil {
		return err
	}

	slf.netWait.Add(1)
	go slf.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (slf *wsServer) serve(ln net.Listener) {
	defer slf.netWait.Done()
	for {
		slf.httpErr = slf.httpSrv.Serve(wsTCPKeepAliveListener{ln.(*net.TCPListener)})
		slf.isShutdown = true
		break
	}

	slf.conns.Range(func(handle interface{}, v interface{}) bool {
		so := operGet(handle.(int32))
		if so.b == resIdle {
			return true
		}

		so.l.Lock()
		if so.b == resIdle || so.b == resOccupy || so.s == nil {
			so.l.Unlock()
			return true
		}

		conn := so.s
		conn.close(nil)
		so.l.Unlock()
		conn.closewait()

		return true
	})
}

func (slf *wsServer) wsAccept(w http.ResponseWriter, r *http.Request) {

	if slf.stat != Connected {
		return
	}

	s, err := slf.waccept.Upgrade(w, r, nil)
	if err != nil {
		//错误日志
		logger.Fatal(slf.operator.ID, "web socket accept fail:%v", err)
		return
	}

	err = slf.accept(s, s.RemoteAddr().Network(), s.RemoteAddr().String())
	if err != nil {
		logger.Fatal(slf.operator.ID, "web socket accept fail:%v", err)
	}
}

func (slf *wsServer) wsMakeConn(handle int32, s interface{}, operator *actor.PID, so *slot, now uint64, stat int32) ISocket {
	conn := &wsConn{}
	conn.h = handle
	conn.s = s
	conn.o = operator
	conn.so = so
	conn.stat = stat
	conn.rv = wsConnRecv
	conn.wr = wsConnWrite
	conn.cls = wsConnClose
	conn.out = make(chan *NetChunk, slf.outChanMax)
	conn.quit = make(chan int)
	conn.i.ReadLastTime = now
	conn.i.WriteLastTime = now

	return conn
}

func (slf *wsServer) httpListen(addr string) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}

	return net.Listen("tcp", addr)
}

func (slf *wsServer) getProto() string {
	return protoWeb
}
func (slf *wsServer) getType() int {
	return CListen
}

func (slf *wsServer) close(lck *mutex.ReSpinLock) {
	slf.httpSrv.Close()
}
