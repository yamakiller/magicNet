package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
)

//NetClientService doc
//@Struct NetClientService @Summary network client service
//@Inherit (service.Service)
//@Inherit (NetClient)
type NetSrvClient struct {
	handler.Service
	NetClient
}

//Initial doc
//@Method Initial @Summary Client service initialization
func (slf *NetSrvClient) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
}

//GetSocket doc
//@Method GetSocket @Summary Returns the client socket
//@Return (int32) socket id
func (slf *NetSrvClient) GetSocket() int32 {
	return 0
}

//SetSocket doc
//@Method SetSocket @Summary Setting the client socket
//@Param (int32) socket id
func (slf *NetSrvClient) SetSocket(sock int32) {

}

//Stoped doc
//@Method Stoped @Summary Notify that the service has been stopped, clear socket relationship / clear ID
//@Param (actor.Context) current service context
//@Param (*actor.PID)    send id
//@Param (interface{})   message
func (slf *NetSrvClient) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf.GetSocket())
	slf.SetSocket(0)
	slf.Service.Stoped(context, sender, message)
}

//Shutdown doc
//@Method Shutdown @Summary Terminate this client service
func (slf *NetSrvClient) Shutdown() {
	slf.Service.Shutdown()
}
