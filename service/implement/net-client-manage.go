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
	Init()
	Size() int
	Grap(h *util.NetHandle) INetClient
	GrapSocket(sock int32) INetClient
	GetHandles() []util.NetHandle
	Erase(h *util.NetHandle)
	Occupy(c INetClient) (*util.NetHandle, error)
	Release(net INetClient)
	Allocer() IAllocer
}

//NetClientManager server client management base class
type NetClientManager struct {
	Malloc IAllocer
}

//Allocer Return to the distributor interface
func (ncm *NetClientManager) Allocer() IAllocer {
	return ncm.Malloc
}
