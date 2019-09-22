package net

import (
	"bytes"

	"github.com/yamakiller/magicNet/engine/actor"
)

type INetConnectionDataStat interface {
	UpdateWrite(tts uint64, bytes uint64)
	UpdateRead(tts uint64, bytes uint64)
	UpdateOnline(tts uint64)
	GetOnline() uint64
}

//INetConnection Network connection interface
type INetConnection interface {
	Name() string
	Socket() int32
	Connection(context actor.Context, addr string, outChanSize int) error
	Write(wrap []byte, length int) error
	GetRecvBufferLimit() int
	GetRecvBuffer() *bytes.Buffer
	GetDataStat() INetConnectionDataStat
	Close()
}
