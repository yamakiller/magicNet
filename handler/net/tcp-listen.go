package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen TCP listener
type TCPListen struct {
	_s int32
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

	slf._s = sock
	return nil
}

// Close Turn off listening
func (slf *TCPListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
