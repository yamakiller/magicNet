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
	slf.s = c
	slf.stat = Connecting
	slf.rv = tcpConnRecv
	slf.wr = tcpConnWrite
	slf.cls = tcpConnClose
	slf.i.ReadLastTime = now
	slf.i.WriteLastTime = now

	return nil
}

func (slf *tcpClient) getProto() string {
	return protoTCP
}

func (slf *tcpClient) getType() int {
	return CClient
}
