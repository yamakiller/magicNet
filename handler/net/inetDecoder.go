package net

import "github.com/yamakiller/magicNet/engine/actor"

//INetDecoder Decoder interface
type INetDecoder func(context actor.Context, param ...interface{}) error
