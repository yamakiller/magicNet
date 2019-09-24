package implement

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/net"
	"github.com/yamakiller/magicNet/timer"
	"github.com/yamakiller/magicNet/util"
)

//INetListenDeleate Network listening commission
type INetListenDeleate interface {
	Handshake(c INetClient) error
	Analysis(context actor.Context, nets *NetListenService, c INetClient) error
	UnOnlineNotification(h util.NetHandle) error
}

// NetListenService Network monitoring service
type NetListenService struct {
	service.Service

	NetListen  net.INetListen
	NetClients INetClientManager
	NetDeleate INetListenDeleate

	Addr                  string //listening address
	CCMax                 int    //Connector pipe buffer to small
	MaxClient             int
	ClientKeep            uint64
	ClientRecvBufferLimit int
}

//Init Initialize the network listening service
func (nets *NetListenService) Init() {
	nets.Service.Init()
	nets.RegisterMethod(&actor.Started{}, nets.Started)
	nets.RegisterMethod(&actor.Stopped{}, nets.Stoped)
	nets.RegisterMethod(&network.NetAccept{}, nets.OnAccept)
	nets.RegisterMethod(&network.NetChunk{}, nets.OnRecv)
	nets.RegisterMethod(&network.NetClose{}, nets.OnClose)
}

func (nets *NetListenService) getDesc() string {
	return fmt.Sprintf("Network Listen [%s] ", nets.NetListen.Name())
}

//Started Turn on network monitoring service
func (nets *NetListenService) Started(context actor.Context, message interface{}) {
	nets.Assignment(context)
	nets.LogInfo("Service Startup %s", nets.Addr)
	err := nets.NetListen.Listen(context, nets.Addr, nets.CCMax)
	if err != nil {
		nets.LogError("Service Startup failed:%s", err.Error())
		return
	}

	nets.Service.Started(context, message)
	nets.LogInfo("%s Service Startup completed", nets.Name())
}

//Stoped Turn off network monitoring service
func (nets *NetListenService) Stoped(context actor.Context, message interface{}) {
	nets.LogInfo("Service Stoping %s", nets.Addr)

	hls := nets.NetClients.GetHandles()
	if hls != nil && len(hls) > 0 {
		for nets.NetClients.Size() > 0 {
			ick := 0
			for i := 0; i < len(hls); i++ {
				c := nets.NetClients.Grap(&hls[i])
				if c == nil {
					continue
				}
				sck := c.GetSocket()
				nets.NetClients.Release(c)
				network.OperClose(sck)
			}

			for {
				time.Sleep(time.Duration(500) * time.Microsecond)
				if nets.NetClients.Size() <= 0 {
					break
				}

				nets.LogInfo("Service The remaining %d connections need to be closed", nets.NetClients.Size())
				ick++
				if ick > 6 {
					break
				}
			}
		}
	}
	nets.NetListen.Close()
	nets.LogInfo("Service Stoped")
}

//OnAccept Receive connection event
func (nets *NetListenService) OnAccept(context actor.Context, message interface{}) {
	accepter := message.(*network.NetAccept)
	if nets.NetClients.Size()+1 > nets.MaxClient {
		nets.LogWarning("OnAccept client fulled:%d", nets.NetClients.Size())
		network.OperClose(accepter.Handle)
		return
	}

	c := nets.NetClients.Allocer().New()
	if c == nil {
		nets.LogError("OnAccept client closed: insufficient memory")
		network.OperClose(accepter.Handle)
		return
	}

	h, err := nets.NetClients.Occupy(c)
	if err != nil {
		nets.LogError("OnAccept client closed: %v, %d-%s:%d",
			err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		nets.NetClients.Allocer().Delete(c)
		network.OperClose(accepter.Handle)
		return
	}

	c.SetSocket(accepter.Handle)
	c.SetAddr(accepter.Addr.String() + strconv.Itoa(accepter.Port))
	c.SetID(h)

	network.OperOpen(accepter.Handle)
	network.OperSetKeep(accepter.Handle, nets.ClientKeep)

	if err = nets.NetDeleate.Handshake(c); err != nil {
		nets.LogError("OnAccept client fail:%s", err)
	}

	c.GetStat().UpdateOnline(timer.Now())

	nets.NetClients.Release(c)

	nets.LogDebug("OnAccept client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

//OnRecv Receiving data events
func (nets *NetListenService) OnRecv(context actor.Context, message interface{}) {
	defer nets.LogDebug("onRecv complete")

	data := message.(*network.NetChunk)
	c := nets.NetClients.GrapSocket(data.Handle)
	if c == nil {
		nets.LogError("OnRecv No target [%d] client object was found", data.Handle)
		return
	}

	defer nets.NetClients.Release(c)

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
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
				nets.LogError("OnRecv error %+v socket %d", err, data.Handle)
				network.OperClose(data.Handle)
				break
			}

			pos += space
			writed += space

			c.GetStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = nets.NetDeleate.Analysis(context, nets, c)
			if err != nil {
				if err == ErrAnalysisSuccess {
					continue
				} else if err != ErrAnalysisProceed {
					nets.LogError("OnRecv error %+v socket %d closing client", err, data.Handle)
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

//OnClose Close connection event
func (nets *NetListenService) OnClose(context actor.Context, message interface{}) {
	closer := message.(*network.NetClose)
	nets.LogDebug("close socket:%d", closer.Handle)
	c := nets.NetClients.GrapSocket(closer.Handle)
	if c == nil {
		nets.LogError("close unfind map-id socket %d", closer.Handle)
		return
	}

	defer nets.NetClients.Release(c)
	hClose := *c.GetID()

	nets.NetClients.Erase(&hClose)

	if err := nets.NetDeleate.UnOnlineNotification(hClose); err != nil {
		nets.LogDebug("closed client Notification %+v", err)
	}

	nets.LogDebug("closed client: %+v", hClose.GetValue())
}

//LogInfo Log information
func (nets *NetListenService) LogInfo(frmt string, args ...interface{}) {
	nets.Service.LogInfo(nets.getDesc()+frmt, args...)
}

//LogError Record error log information
func (nets *NetListenService) LogError(frmt string, args ...interface{}) {
	nets.Service.LogError(nets.getDesc()+frmt, args...)
}

//LogDebug Record debug log information
func (nets *NetListenService) LogDebug(frmt string, args ...interface{}) {
	nets.Service.LogDebug(nets.getDesc()+frmt, args...)
}

//LogTrace Record trace log information
func (nets *NetListenService) LogTrace(frmt string, args ...interface{}) {
	nets.Service.LogTrace(nets.getDesc()+frmt, args...)
}

//LogWarning Record warning log information
func (nets *NetListenService) LogWarning(frmt string, args ...interface{}) {
	nets.Service.LogWarning(nets.getDesc()+frmt, args...)
}
