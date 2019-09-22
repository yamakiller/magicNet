package net

//TODO:将删除
/*import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

// TCPService :  TCP network listening service
type TCPService struct {
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

// Init TCP network service initialization
func (ts *TCPService) Init() {
	ts.Service.Init()
	ts.RegisterMethod(&actor.Started{}, ts.Started)
	ts.RegisterMethod(&actor.Stopped{}, ts.Stoped)
	ts.RegisterMethod(&network.NetAccept{}, ts.OnAccept)
	ts.RegisterMethod(&network.NetChunk{}, ts.OnRecv)
	ts.RegisterMethod(&network.NetClose{}, ts.OnClose)
}

// Started TCP network service is enabled
func (ts *TCPService) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network Listen [TCP/IP] Service Startup %s", ts.Addr)
	sock, err := network.OperTCPListen(context.Self(), ts.Addr, ts.CCMax)
	if err != nil {
		logger.Error(context.Self().GetID(), "Network Listen [TCP/IP] Service Startup failed:%s", err.Error())
		return
	}

	ts.sock = sock
	ts.Service.Started(context, message)
	logger.Info(context.Self().GetID(), "Network Listen [TCP/IP] Service Startup completed")
}

// Stoped TCP network service stops
func (ts *TCPService) Stoped(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network Listen [TCP/IP] Service Stoping %s", ts.Addr)
	if ts.sock != 0 {
		network.OperClose(ts.sock)
		ts.sock = 0
	}
	logger.Info(context.Self().GetID(), "Network Listen [TCP/IP] Service Stoped")
}

// Shutdown TCP network service termination
func (ts *TCPService) Shutdown() {
	ts.Service.Shutdown()
}*/
