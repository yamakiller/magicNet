package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service"
)

//NetClientService desc
//@struct NetClientService desc: network client service
type NetClientService struct {
	service.Service
	NetClient
}

//Initial Client service initialization
func (slf *NetClientService) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
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
func (slf *NetClientService) GetSocket() int32 {
	return 0
}

//SetSocket Setting the client socket
func (slf *NetClientService) SetSocket(sock int32) {

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
func (slf *NetClientService) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf.GetSocket())
	slf.SetSocket(0)
	slf.Service.Stoped(context, sender, message)
}

//Shutdown Terminate this client service
func (slf *NetClientService) Shutdown() {
	slf.Service.Shutdown()
}
