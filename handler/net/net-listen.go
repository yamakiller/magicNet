package net

import "github.com/yamakiller/magicNet/engine/actor"

//INetListen Network listening interface
type INetListen interface {
	Name() string
	Listen(context actor.Context, addr string, ccmax int) error
	Close()
}
