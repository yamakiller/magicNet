package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//WSSListen doc
//@Summary WebSocket Listen
//@Struct WSSListen
//@Member (int32) service socket id
type WSSListen struct {
	_s int32
}

//GetSocket doc
//@Summary Return socket id
//@Return int32 socket
func (slf *WSSListen) GetSocket() int32 {
	return slf._s
}

//Listen doc
//@Summary Start listening
//@Method Listen
//@Param (actor.Context) Service context
//@Param (string) listen address
//@Param (int)    Revice Data Chan size
//@Return (error) listen fail return error
func (slf *WSSListen) Listen(context actor.Context,
	addr string,
	ccmax int) error {

	sock, err := network.OperWSListen(context.Self(), addr, ccmax)
	if err != nil {
		return err
	}

	slf._s = sock
	return nil
}

//Close doc
//@Summary Turn off listening
//@Method Close
func (slf *WSSListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}

//ToString doc
//@Summary Features
//@Method ToString
//@Param (string)
func (slf *WSSListen) ToString() string {
	return "WebSocket"
}
