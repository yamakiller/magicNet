package net

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen TCP listener
type TCPListen struct {
	s int32
}

//Name Features
func (tpl *TCPListen) Name() string {
	return "TCP/IP"
}

// Listen Start listening
func (tpl *TCPListen) Listen(context actor.Context, addr string, ccmax int) error {
	sock, err := network.OperTCPListen(context.Self(), addr, ccmax)
	if err != nil {
		return err
	}

	tpl.s = sock
	return nil
}

// Close Turn off listening
func (tpl *TCPListen) Close() {
	if tpl.s != 0 {
		network.OperClose(tpl.s)
		tpl.s = 0
	}
}
