package implement

import (
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service/net"
	"github.com/yamakiller/magicNet/network"
)

type NetHandshakeReply func(sock int32)
type NetAnalysis func(c INetClient) error

// NetService
type NetService struct {
	service.Service

	NetListen  			 net.INetListen
	NetClients 			 INetClientManager
	OnHandshakeReply NetHandshakeReply
	OnAnalysis       NetAnalysis

	Addr  		 string //listening address
	CCMax 	   int    //Connector pipe buffer to small
	MaxClient  int
	ClientKeep uint64
	ClientRecvBufferLimit int
}


func (nets *NetService) Init() {
	nets.Service.Init()
	nets.RegisterMethod(&actor.Started{}, nets.Started)
	nets.RegisterMethod(&actor.Stopped{}, nets.Stoped)
	nets.RegisterMethod(&network.NetAccept{}, nets.OnAccept)
	nets.RegisterMethod(&network.NetChunk{}, nets.OnRecv)
	nets.RegisterMethod(&network.NetClose{}, nets.OnClose)
}

func (nets *NetService) getDesc() string {
	return fmt.Sprintf("Network Listen [%s]", nets.NetListen.Name())
}

func (nets *NetService) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "%s Service Startup %s", nets.getDesc(), nets.Addr)
	err := nets.NetListen.Listen(context, nets.Addr, nets.CCMax)
	if err != nil {
		logger.Error(context.Self().GetID(), "%s Service Startup failed:%s",  nets.getDesc(), err.Error())
		return
	}

	nets.Service.Started(context, message)
	logger.Info(context.Self().GetID(), "%s Service Startup completed", nets.getDesc(), nets.Name())
}

func (nets *NetService) Stoped(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "%s Service Stoping %s", nets.getDesc(), nets.Addr)
	//关闭所有客户端连接
	nets.NetListen.Close()
	logger.Info(context.Self().GetID(), "%s Service Stoped", nets.getDesc())
}

func (nets *NetService) OnAccept(context actor.Context, message interface{}) {
	accepter := message.(*network.NetAccept)
	if nets.NetClients.Size() + 1 > nets.MaxClient {
		logger.Warning(context.Self().GetID(), "%s OnAccept client fulled:%d",  nets.getDesc(), nets.NetClients.Size())
		network.OperClose(accepter.Handle)
		return
	}

	c := nets.NetClients.Allocer().New()
	if c == nil {
		logger.Error(context.Self().GetID(), "%s OnAccept client closed: insufficient memory", nets.getDesc())
		network.OperClose(accepter.Handle)
		return
	}

	h , err := nets.NetClients.Occupy()
	if err != nil {
		logger.Error(context.Self().GetID(), "%s OnAccept client closed: %v, %d-%s:%d", nets.getDesc(), err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		nets.NetClients.Allocer().Delete(c)
		network.OperClose(accepter.Handle)
		return
	}

	c.SetSocket(accepter.Handle)
	c.SetAddr(accepter.Addr.String())
	c.SetPort(accepter.Port)
	c.SetID(h)

	network.OperOpen(accepter.Handle)
	network.OperSetKeep(accepter.Handle, nets.ClientKeep)

	nets.OnHandshakeReply(accepter.Handle)

	c.GetStat().UpdateOnline(timer.Now())

	nets.NetClients.Release(c)

	logger.Debug(context.Self().GetID(), "%s OnAccept client %d-%s:%d", nets.getDesc(), accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (nets *NetService) OnRecv(context actor.Context, message interface{}) {
	defer logger.Debug(self.Self().GetID(), "%s onRecv complete", nets.getDesc())
	
	data := message.(*network.NetChunk)
	c := nets.NetClients.GrapSocket(data.Handle)
	if c == nil {
		logger.Error(context.Self().GetID(), "%s OnRecv No target [%d] client object was found", nets.getDesc(), data.Handle)
		return
	}

	defer nets.NetClients.Release(c)

	var (
		space  int
		writed int
		wby    int
		pos    int

		err    error
	)

	for {
			space = nets.ClientRecvBufferLimit - c.GetRecvBuffer().Len()
			wby = len(data.Data) - writed
			if space > 0 && wby > 0 {
				if space > wby {
					space = wby
				}

				_, err = c.GetRecvBuffer().Write(data.Data[pos : pos+space])
				if err != nil {
					logger.Error(context.Self().GetID(), "%s OnRecv error %+v socket %d", nets.getDesc(), err, data.Handle)
					network.OperClose(data.Handle)
					break
				}

				pos += space
				writed += space

				c.GetStat().UpdateRead(timer.Now(), uint64(space))
			}

			for {
				// Decomposition of Packets
				err = nets.OnAnalysis(c)
				if err != nil {
					if err == AnalysisSuccess {
						continue
					}  else if err != AnalysisProceed {
						logger.Error(context.Self().GetID(), "%s OnRecv error %+v socket %d closing client", nets.getDesc(), err, data.Handle)
						network.OperClose(data.Handle)
						return
					}
				}

				if writed >= len(data.Data) {
					return
				}

				break
			}
	}
}

func (nets *NetService) OnClose(context actor.Context, message interface{}) {

}
