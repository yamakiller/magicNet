package client

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
)

type NetSSrvCleint struct {
	handler.Service
	NetSrvClient
}

//Initial doc
//@Summary Server Client service initialization
//@Method Initial
func (slf *NetSSrvCleint) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
}

//Stoped doc
//@Summary Notify that the service has been stopped, clear socket relationship / clear ID
//@Method Stoped
//@Param (actor.Context) current service context
//@Param (*actor.PID)    send id
//@Param (interface{})   message
func (slf *NetSSrvCleint) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf.GetSocket())
	slf.WithSocket(net.INVALIDSOCKET)
	slf.Service.Stoped(context, sender, message)
}

//Shutdown doc
//@Summary Terminate this client service
//@Method Shutdown
func (slf *NetSSrvCleint) Shutdown() {
	slf.Service.Shutdown()
}
