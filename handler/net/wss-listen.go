package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//WSSListen desc
//@struct WSSListen desc: WebSocket Listen
//@member (int32) service socket id
type WSSListen struct {
	_s int32
}

//Name desc
//@method Name desc: Features
//@param (string)
func (slf *WSSListen) Name() string {
	return "WebSocket"
}

//Listen desc
//@method Listen desc: Start listening
//@param (actor.Context) Service context
//@param (string) listen address
//@param (int)    Revice Data Chan size
//@return (error) listen fail return error
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

//Close desc
//@method Close desc: Turn off listening
func (slf *WSSListen) Close() {
	if net.InvalidSocket(slf._s) {
		network.OperClose(slf._s)
		slf._s = net.INVALIDSOCKET
	}
}
