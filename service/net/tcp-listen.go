package net

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

type TCPListen struct {
	s int32
}

func (tpl *TCPListen) Name() string {
	return "TCP/IP"
}

func (tpl *TCPListen) Listen(context actor.Context, addr string, ccmax int) error {
	sock, err := network.OperTCPListen(context.Self(), addr, ccmax)
	if err != nil {
		return err
	}

	tpl.s = sock
	return nil
}

func (tpl *TCPListen) Close() {
	if tpl.s != 0 {
		network.OperClose(tpl.s)
		tpl.s = 0
	}
}
