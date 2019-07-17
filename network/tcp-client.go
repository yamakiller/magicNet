package network

import (
	"magicNet/engine/actor"
	"magicNet/timer"
	"net"
)

type tcpClient struct {
	sConn
}

func (tpc *tcpClient) connect(operator *actor.PID, addr string) error {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	now := timer.Now()
	tpc.s = c
	tpc.stat = Connecting
	tpc.rv = tcpConnRecv
	tpc.wr = tcpConnWrite
	tpc.cls = tcpConnClose
	tpc.i.ReadLastTime = now
	tpc.i.WriteLastTime = now

	return nil
}
