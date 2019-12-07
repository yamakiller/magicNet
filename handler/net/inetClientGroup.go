package net

//IAllocer Distributor interface
type IAllocer interface {
	New() INetClient
	Delete(p INetClient)
}

//INetClientGroup Network client management interface
type INetClientGroup interface {
	Initial()
	Size() int
	Cap() int
	Grap(h uint64) INetClient
	GrapSocket(sock int32) INetClient
	GetHandles() []uint64
	Erase(h uint64)
	Occupy(c INetClient) (uint64, error)
	Release(net INetClient)
	Allocer() IAllocer
}

//NetClientManager server client management base class
/*type NetClientGroup struct {
	Malloc IAllocer
}

//Allocer Return to the distributor interface
func (slf *NetClientGroup) Allocer() IAllocer {
	return slf.Malloc
}*/
