package netservice

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

//WebConnection connection service
type WebConnection struct {
	service.Service

	OnRecv  service.MethodFunc //NetChunk
	OnClose service.MethodFunc //NetClose
}

// Init TCP network service initialization
func (wc *WebConnection) Init() {
	wc.Service.Init()
	wc.RegisterMethod(&actor.Started{}, wc.Started)
	wc.RegisterMethod(&actor.Stopped{}, wc.Stoped)
	wc.RegisterMethod(&network.NetChunk{}, wc.OnRecv)
	wc.RegisterMethod(&network.NetClose{}, wc.OnClose)
}

//Started Start Web Socket connection service
func (wc *WebConnection) Started(context actor.Context, message interface{}) {
	wc.Service.Started(context, message)
}

// Stoped Web Socket network service stops
func (wc *WebConnection) Stoped(context actor.Context, message interface{}) {
	wc.Service.Stoped(context, message)
}

// Shutdown Web Socket network service termination
func (wc *WebConnection) Shutdown() {
	wc.Service.Shutdown()
}
