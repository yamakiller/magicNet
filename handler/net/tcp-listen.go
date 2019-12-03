package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen desc:
//@Struct TCPListen desc: TCP listener
//@Member (int32) TCP listen service socket id
type TCPListen struct {
	_s int32
}

//Name desc
//@Method Name desc: Features
//@Return (string)
func (slf *TCPListen) Name() string {
	return "TCP/IP"
}

//Listen desc
//@Method Listen desc: Start listening
//@Param (actor.Context) Service context
//@Param (string) listen address
//@Param (int)    Revice Data Chan size
//@Return (error) listen fail return error
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
//@Method Close desc: Turn off listening
func (slf *TCPListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
