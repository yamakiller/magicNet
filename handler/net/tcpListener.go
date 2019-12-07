package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPListen @Summary
//@Summary TCP listener
//@Struct TCPListen
//@Member (int32) TCP listen service socket id
type TCPListen struct {
	_s int32
}

//GetSocket doc
//@Summary Return socket id
//@Return int32 socket
func (slf *TCPListen) GetSocket() int32 {
	return slf._s
}

//Listen doc
//@Summary Start listening
//@Method Listen
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
//@Summary Turn off listening
//@Method Close
func (slf *TCPListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}

//ToString doc
//@Summary Features
//@Method ToString
//@Return (string)
func (slf *TCPListen) ToString() string {
	return "TCP/IP"
}
