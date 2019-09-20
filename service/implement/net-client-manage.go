package implement

import (
  "github.com/yamakiller/magicNet/util"
)

type IAllocer interface {
  New() INetClient
  Delete(p INetClient)
}

type INetClientManager interface {
  Spawn()
  Size() int
  Grap(h *util.NetHandle) INetClient
  GrapSocket(sock int32) INetClient
  Erase(h *util.NetHandle)
  Occupy(c INetClient) (*util.NetHandle ,error)
  Release(net INetClient)
  Allocer() IAllocer
}
