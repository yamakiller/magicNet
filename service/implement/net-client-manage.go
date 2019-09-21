package implement

import (
	"github.com/yamakiller/magicNet/util"
)

//IAllocer Distributor interface
type IAllocer interface {
	New() INetClient
	Delete(p INetClient)
}

//INetClientManager Network client management interface
type INetClientManager interface {
	Size() int
	Grap(h *util.NetHandle) INetClient
	GrapSocket(sock int32) INetClient
	GetHandles() []util.NetHandle
	Erase(h *util.NetHandle)
	Occupy(c INetClient) (*util.NetHandle, error)
	Release(net INetClient)
	Allocer() IAllocer
}

//NetClientManager
type NetClientManager struct {
	Malloc IAllocer
}

func (ncm *NetClientManager) Allocer() IAllocer {
	return ncm.Malloc
}
