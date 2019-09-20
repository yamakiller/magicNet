package net

type INetListen interface {
   Name() string
   Listen(context actor.Context, addr string, ccmax int) error
   Close()
}
