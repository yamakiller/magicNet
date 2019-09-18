package netservice

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

//TCPConnection connection service
type TCPConnection struct {
	service.Service

	OnRecv  service.MethodFunc //NetChunk
	OnClose service.MethodFunc //NetClose
}

// Init TCP network service initialization
func (tc *TCPConnection) Init() {
	tc.Service.Init()
	tc.RegisterMethod(&actor.Started{}, tc.Started)
	tc.RegisterMethod(&actor.Stopped{}, tc.Stoped)
	tc.RegisterMethod(&network.NetChunk{}, tc.OnRecv)
	tc.RegisterMethod(&network.NetClose{}, tc.OnClose)
}

//Started Start tcp connection service
func (tc *TCPConnection) Started(context actor.Context, message interface{}) {
	tc.Service.Started(context, message)
}

// Stoped TCP network service stops
func (tc *TCPConnection) Stoped(context actor.Context, message interface{}) {
	tc.Service.Stoped(context, message)
}

// Shutdown TCP network service termination
func (tc *TCPConnection) Shutdown() {
	tc.Service.Shutdown()
}
