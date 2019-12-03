package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
)

//NetClientService desc
//@Struct NetClientService desc: network client service
//@Inherit (service.Service)
//@Inherit (NetClient)
type NetClientService struct {
	handler.Service
	NetClient
}

//Initial desc
//@Method Initial desc: Client service initialization
func (slf *NetClientService) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
}

//GetSocket desc
//@Method GetSocket desc: Returns the client socket
//@Return (int32) socket id
func (slf *NetClientService) GetSocket() int32 {
	return 0
}

//SetSocket desc
//@Method SetSocket desc: Setting the client socket
//@Param (int32) socket id
func (slf *NetClientService) SetSocket(sock int32) {

}

//Stoped desc
//@Method Stoped desc: Notify that the service has been stopped, clear socket relationship / clear ID
//@Param (actor.Context) current service context
//@Param (*actor.PID)    send id
//@Param (interface{})   message
func (slf *NetClientService) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf.GetSocket())
	slf.SetSocket(0)
	slf.Service.Stoped(context, sender, message)
}

//Shutdown desc
//@Method Shutdown desc: Terminate this client service
func (slf *NetClientService) Shutdown() {
	slf.Service.Shutdown()
}
