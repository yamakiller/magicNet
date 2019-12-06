package net

import (
	"github.com/yamakiller/magicNet/engine/actor"
)

//INetConnectionDataStat Connector data status maintenance interface
type INetConnectionDataStat interface {
	UpdateWrite(tts uint64, bytes uint64)
	UpdateRead(tts uint64, bytes uint64)
	UpdateOnline(tts uint64)
	GetOnline() uint64
}

//INetConnection Network connection interface
type INetConnection interface {
	Connection(context actor.Context,
		addr string, /*Connection address*/
		outChanSize int /*Receive pipe buffer size*/) error
	Write(wrap []byte, length int) error

	GetSocket() int32

	Close()

	// GetReceiveBufferLimit() int
	// GetReceiveBuffer() *bytes.Buffer
}
