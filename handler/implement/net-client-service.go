package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
)

//NetClientService desc
//@struct NetClientService desc: network client service
//@inherit (service.Service)
//@inherit (NetClient)
type NetClientService struct {
	handler.Service
	NetClient
}

//Initial desc
//@method Initial desc: Client service initialization
func (slf *NetClientService) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
}

//GetSocket desc
//@method GetSocket desc: Returns the client socket
//@return (int32) socket id
func (slf *NetClientService) GetSocket() int32 {
	return 0
}

//SetSocket desc
//@method SetSocket desc: Setting the client socket
//@param (int32) socket id
func (slf *NetClientService) SetSocket(sock int32) {

}

//Stoped desc
//@method Stoped desc: Notify that the service has been stopped, clear socket relationship / clear ID
//@param (actor.Context) current service context
//@param (*actor.PID)    send id
//@param (interface{})   message
func (slf *NetClientService) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf.GetSocket())
	slf.SetSocket(0)
	slf.Service.Stoped(context, sender, message)
}

//Shutdown desc
//@method Shutdown desc: Terminate this client service
func (slf *NetClientService) Shutdown() {
	slf.Service.Shutdown()
}
