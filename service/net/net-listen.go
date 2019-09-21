package net

import "github.com/yamakiller/magicNet/engine/actor"

type INetListen interface {
	Name() string
	Listen(context actor.Context, addr string, ccmax int) error
	Close()
}
