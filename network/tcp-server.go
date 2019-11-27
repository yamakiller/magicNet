package network

import (
	"net"
	"time"

	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
)

type tcpServer struct {
	sServer
	_s *net.TCPListener
}

func (slf *tcpServer) listen(operator *actor.PID, addr string) error {

	tcpAddr, aderr := net.ResolveTCPAddr("tcp", addr)
	if aderr != nil {
		return aderr
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	slf._s = ln
	slf._maker = slf.makeConn

	slf._netWait.Add(1)
	go slf.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (slf *tcpServer) serve(ln net.Listener) {
	defer slf._netWait.Done()
	for {

		s, err := slf._s.AcceptTCP()
		if err != nil {
			logger.Error(slf._operator.ID, "socket accept fail:%s", err.Error())
			slf._isShutdown = true
			break
		}
		s.SetNoDelay(true)

		err = slf.accept(s, s.RemoteAddr().Network(), s.RemoteAddr().String())
		if err != nil {
			logger.Fatal(slf._operator.ID, "socket accept fail:%v", err)
		}
	}

	//------------------关闭所有连接-----------------------------
	slf._conns.Range(func(handle interface{}, v interface{}) bool {
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

func (slf *tcpServer) keeploop() {

}

func (slf *tcpServer) makeConn(handle int32,
	s interface{},
	operator *actor.PID,
	so *slot,
	now uint64,
	stat int32) ISocket {

	conn := &tcpConn{}
	conn._h = handle
	conn._s = s
	conn._o = operator
	conn._so = so
	conn._stat = stat
	conn._rv = tcpConnRecv
	conn._wr = tcpConnWrite
	conn._cls = tcpConnClose
	conn._out = make(chan *NetChunk, slf._outChanMax)
	conn._quit = make(chan int)
	conn._i.RecvLastTime = now
	conn._i.WriteLastTime = now
	conn._w.Add(2)
	return conn
}

func (slf *tcpServer) getProto() string {
	return protoTCP
}
func (slf *tcpServer) getType() int {
	return CListen
}

func (slf *tcpServer) close(lck *mutex.ReSpinLock) {
	slf._s.Close()
}
