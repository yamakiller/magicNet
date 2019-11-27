package network

import (
	"net"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"
)

type tcpClient struct {
	sConn
}

func (slf *tcpClient) connect(operator *actor.PID, addr string) error {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	now := timer.Now()
	slf._s = c
	slf._stat = Connecting
	slf._rv = tcpConnRecv
	slf._wr = tcpConnWrite
	slf._cls = tcpConnClose
	slf._i.RecvLastTime = now
	slf._i.WriteLastTime = now

	return nil
}

func (slf *tcpClient) getProto() string {
	return protoTCP
}

func (slf *tcpClient) getType() int {
	return CClient
}
