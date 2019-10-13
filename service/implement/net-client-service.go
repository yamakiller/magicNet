package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service"
)

//NetClientService 网络客户端服务层
type NetClientService struct {
	service.Service
	NetClient
}

//Init Client service initialization
func (ncs *NetClientService) Init() {
	ncs.Service.Init()
	ncs.RegisterMethod(&actor.Stopped{}, ncs.Stoped)
}

/*//SetID Setting the client ID
func (ncs *NetClientService) SetID(h uint64) {
	ncs.handle.SetValue(h)
}

//GetID Returns the client ID
func (ncs *NetClientService) GetID() uint64 {
	//return ncs.handle.GetValue()
}*/

//GetSocket Returns the client socket
func (ncs *NetClientService) GetSocket() int32 {
	return 0
}

//SetSocket Setting the client socket
func (ncs *NetClientService) SetSocket(sock int32) {

}

//GetAuth return to certification time
/*func (ncs *NetClientService) GetAuth() uint64 {
	return 0
}

//SetAuth Setting the time for authentication
func (ncs *NetClientService) SetAuth(v uint64) {
}*/

//GetKeyPair Return key object
/*func (ncs *NetClientService) GetKeyPair() interface{} {
	return nil
}

//BuildKeyPair Build key pair
func (ncs *NetClientService) BuildKeyPair() {

}

//GetKeyPublic Return key publicly available information
func (ncs *NetClientService) GetKeyPublic() string {
	return ""
}*/

//Stoped 通知已经停止服务，清除套接字关系/清除ID
func (ncs *NetClientService) Stoped(context actor.Context, message interface{}) {
	ncs.LogDebug("Stoped: Socket-%d", ncs.GetSocket())
	ncs.SetSocket(0)
	ncs.Service.Stoped(context, message)
}

//Shutdown Terminate this client service
func (ncs *NetClientService) Shutdown() {
	ncs.Service.Shutdown()
}
