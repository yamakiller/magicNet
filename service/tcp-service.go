package service

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
)

// TCPService : 网络服务
type TCPService struct {
	Service
	sock int32
	//
	Addr  string //监听地址[IP:Port]
	CCMax int    //连接者管道缓冲区到小
	//
	OnAccept MethodFunc //NetAccept
	OnRecv   MethodFunc //NetChunk
	OnClose  MethodFunc //NetClose
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
	logger.Info(context.Self().GetID(), "Network[TCP/IP] Service Start %s", ts.Addr)
	sock, err := network.OperTCPListen(context.Self(), ts.Addr, ts.CCMax)
	if err != nil {
		logger.Error(context.Self().GetID(), "Network[TCP/IP] Service Start Fail:%s", err.Error())
		return
	}

	ts.sock = sock
	ts.Service.Started(context, message)
	logger.Info(context.Self().GetID(), "Network[TCP/IP] Service Success")
}

// Stoped TCP network service stops
func (ts *TCPService) Stoped(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network[TCP/IP] Service Stoping %s", ts.Addr)
	if ts.sock != 0 {
		network.OperClose(ts.sock)
		ts.sock = 0
	}
	logger.Info(context.Self().GetID(), "Network[TCP/IP] Service Stoped")
}

// Shutdown TCP network service termination
func (ts *TCPService) Shutdown() {
	ts.Service.Shutdown()
}
