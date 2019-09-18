package netservice

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

// WebSocketService :  Web socket network listening service
type WebSocketService struct {
	service.Service
	sock int32
	//
	Addr  string //listening address  [IP:Port]
	CCMax int    //Connector pipe buffer to small
	//
	OnAccept service.MethodFunc //NetAccept
	OnRecv   service.MethodFunc //NetChunk
	OnClose  service.MethodFunc //NetClose
}

// Init Web Socket network service initialization
func (wss *WebSocketService) Init() {
	wss.Service.Init()
	wss.RegisterMethod(&actor.Started{}, wss.Started)
	wss.RegisterMethod(&actor.Stopped{}, wss.Stoped)
	wss.RegisterMethod(&network.NetAccept{}, wss.OnAccept)
	wss.RegisterMethod(&network.NetChunk{}, wss.OnRecv)
	wss.RegisterMethod(&network.NetClose{}, wss.OnClose)
}

// Started Web Socket network service is enabled
func (wss *WebSocketService) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network Listen [Web/Socket] Service Startup %s", wss.Addr)
	sock, err := network.OperWSListen(context.Self(), wss.Addr, wss.CCMax)
	if err != nil {
		logger.Error(context.Self().GetID(), "Network Listen [Web/Socket] Service Startup failed:%s", err.Error())
		return
	}

	wss.sock = sock
	wss.Service.Started(context, message)
	logger.Info(context.Self().GetID(), "Network Listen [Web/Socket] Service Startup completed")
}

// Stoped Web Socket network service stops
func (wss *WebSocketService) Stoped(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network Listen [Web/Socket] Service Stoping %s", wss.Addr)
	if wss.sock != 0 {
		network.OperClose(wss.sock)
		wss.sock = 0
	}
	logger.Info(context.Self().GetID(), "Network Listen [Web/Socket] Service Stoped")
}

// Shutdown Web Socket network service termination
func (wss *WebSocketService) Shutdown() {
	wss.Service.Shutdown()
}
