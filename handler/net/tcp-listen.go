package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen @Summary
//@Struct TCPListen @Summary TCP listener
//@Member (int32) TCP listen service socket id
type TCPListen struct {
	_s int32
}

//Name doc
//@Method Name @Summary Features
//@Return (string)
func (slf *TCPListen) Name() string {
	return "TCP/IP"
}

//Listen doc
//@Method Listen @Summary Start listening
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

//Close doc
//@Method Close @Summary Turn off listening
func (slf *TCPListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
