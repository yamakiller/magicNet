package net

import "github.com/yamakiller/magicNet/engine/actor"

//INetListener doc
//@Summary Network listening interface
//@Interface INetListener
//@Method ToString() string
//@Method Listen
//@Method Close
type INetListener interface {
	ToString() string
	GetSocket() int32
	Listen(context actor.Context, addr string, ccmax int) error
	Close()
}
