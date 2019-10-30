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
func (slf *TCPListen) Name() string {
	return "TCP/IP"
}

// Listen Start listening
func (slf *TCPListen) Listen(context actor.Context,
	addr string,
	ccmax int) error {

	sock, err := network.OperTCPListen(context.Self(), addr, ccmax)
	if err != nil {
		return err
	}

	slf.s = sock
	return nil
}

// Close Turn off listening
func (slf *TCPListen) Close() {
	if slf.s != 0 {
		network.OperClose(slf.s)
		slf.s = 0
	}
}
