package net

import (
	"github.com/yamakiller/magicNet/engine/actor"
)

//INetConnection Network connection interface
type INetConnection interface {
	ToString() string
	Connection(context actor.Context,
		addr string, /*Connection address*/
		outChanSize int /*Receive pipe buffer size*/) error
	Write(wrap []byte, length int) error

	WithSocket(int32)
	GetSocket() int32

	Close()
}
