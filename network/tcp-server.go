package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"net"
	"time"
)

type tcpServer struct {
	sServer
	s net.Listener
}

func (tps *tcpServer) listen(operator *actor.PID, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tps.s = ln
	tps.maker = tps.makeConn

	tps.netWait.Add(2)
	go tps.serve(ln)
	go tps.keeploop()

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (tps *tcpServer) serve(ln net.Listener) {
	defer tps.netWait.Done()
	for {

		s, err := tps.s.Accept()
		if err != nil {
			logger.Error(tps.operator.ID, "socket accept fail:%s", err.Error())
			break
		}

		err = tps.accept(s, s.RemoteAddr().Network(), s.RemoteAddr().String())
		if err != nil {
			logger.Fatal(tps.operator.ID, "socket accept fail:%v", err)
		}
	}
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
	conn.i.ReadLastTime = now
	conn.i.WriteLastTime = now
	return conn
}

func (tps *tcpServer) getProto() string {
	return ProtoTCP
}
func (tps *tcpServer) getType() int {
	return CListen
}

func (tps *tcpServer) close(lck *util.ReSpinLock) {
	tps.s.Close()
}
