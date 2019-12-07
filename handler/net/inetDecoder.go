package net

import "github.com/yamakiller/magicNet/engine/actor"

type INetDecoder func(context actor.Context, param ...interface{}) error
