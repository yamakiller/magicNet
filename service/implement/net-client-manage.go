package implement

//IAllocer Distributor interface
type IAllocer interface {
	New() INetClient
	Delete(p INetClient)
}

//INetClientManager Network client management interface
type INetClientManager interface {
	Init()
	Size() int
	Grap(h uint64) INetClient
	GrapSocket(sock int32) INetClient
	GetHandles() []uint64
	Erase(h uint64)
	Occupy(c INetClient) (uint64, error)
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
