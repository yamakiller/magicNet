package network

import (
	"net"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/util"
)

type tcpServer struct {
	sServer
	s *net.TCPListener
}

func (tps *tcpServer) listen(operator *actor.PID, addr string) error {

	tcpAddr, aderr := net.ResolveTCPAddr("tcp", addr)
	if aderr != nil {
		return aderr
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	tps.s = ln
	tps.maker = tps.makeConn

	tps.netWait.Add(1)
	go tps.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (tps *tcpServer) serve(ln net.Listener) {
	defer tps.netWait.Done()
	for {

		s, err := tps.s.AcceptTCP()
		if err != nil {
			logger.Error(tps.operator.ID, "socket accept fail:%s", err.Error())
			tps.isShutdown = true
			break
		}
		s.SetNoDelay(true)

		err = tps.accept(s, s.RemoteAddr().Network(), s.RemoteAddr().String())
		if err != nil {
			logger.Fatal(tps.operator.ID, "socket accept fail:%v", err)
		}
	}

	//------------------关闭所有连接-----------------------------
	tps.conns.Range(func(handle interface{}, v interface{}) bool {
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

func (tps *tcpServer) keeploop() {

}

func (tps *tcpServer) makeConn(handle int32, s interface{}, operator *actor.PID, so *slot, now uint64, stat int32) ISocket {
	conn := &tcpConn{}
	conn.h = handle
	conn.s = s
	conn.o = operator
	conn.so = so
	conn.stat = stat
	conn.rv = tcpConnRecv
	conn.wr = tcpConnWrite
	conn.cls = tcpConnClose
	conn.out = make(chan *NetChunk, tps.outChanMax)
	conn.quit = make(chan int)
	conn.i.ReadLastTime = now
	conn.i.WriteLastTime = now
	conn.w.Add(2)
	return conn
}

func (tps *tcpServer) getProto() string {
	return protoTCP
}
func (tps *tcpServer) getType() int {
	return CListen
}

func (tps *tcpServer) close(lck *util.ReSpinLock) {
	tps.s.Close()
}
