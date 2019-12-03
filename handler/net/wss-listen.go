package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//WSSListen doc
//@Struct WSSListen @Summary WebSocket Listen
//@Member (int32) service socket id
type WSSListen struct {
	_s int32
}

//Name doc
//@Method Name @Summary Features
//@Param (string)
func (slf *WSSListen) Name() string {
	return "WebSocket"
}

//Listen doc
//@Method Listen @Summary Start listening
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
//@Method Close @Summary Turn off listening
func (slf *WSSListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
