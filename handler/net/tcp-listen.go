package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen desc:
//@struct TCPListen desc: TCP listener
//@member (int32) TCP listen service socket id
type TCPListen struct {
	_s int32
}

//Name desc
//@method Name desc: Features
//@return (string)
func (slf *TCPListen) Name() string {
	return "TCP/IP"
}

//Listen desc
//@method Listen desc: Start listening
//@param (actor.Context) Service context
//@param (string) listen address
//@param (int)    Revice Data Chan size
//@return (error) listen fail return error
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

//Close desc
//@method Close desc: Turn off listening
func (slf *TCPListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
