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
)

//INetListenDeleate Network listening commission
type INetListenDeleate interface {
	Handshake(c INetClient) error
	Analysis(context actor.Context, nets *NetListenService, c INetClient) error
	UnOnlineNotification(h uint64) error
}

// NetListenService Network monitoring service
type NetListenService struct {
	service.Service

	NetListen  net.INetListen
	NetClients INetClientManager
	NetDeleate INetListenDeleate
	NetMethod  NetMethodDispatch

	Addr       string //listening address
	CCMax      int    //Connector pipe buffer to small
	MaxClient  int
	ClientKeep uint64
}

//Initial Initialize the network listening service
func (slf *NetListenService) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&network.NetAccept{}, slf.OnAccept)
	slf.RegisterMethod(&network.NetChunk{}, slf.OnRecv)
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
}

func (slf *NetListenService) getDesc() string {
	return fmt.Sprintf("Network Listen [%s] ", slf.NetListen.Name())
}

//Started Turn on network monitoring service
func (slf *NetListenService) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf.WithPID(context)
	slf.LogInfo("Service Startup %s", slf.Addr)
	err := slf.NetListen.Listen(context, slf.Addr, slf.CCMax)
	if err != nil {
		slf.LogError("Service Startup failed:%s", err.Error())
		return
	}

	slf.Service.Started(context, sender, message)
	slf.LogInfo("%s Service Startup completed", slf.Name())
}

//Stopping Turn off network monitoring service
func (slf *NetListenService) Stopping(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	slf.LogInfo("Service Stoping %s", slf.Addr)

	h := NetHandle{}
	hls := slf.NetClients.GetHandles()
	if hls != nil && len(hls) > 0 {
		for slf.NetClients.Size() > 0 {
			ick := 0
			for i := 0; i < len(hls); i++ {
				h.SetValue(hls[i])
				c := slf.NetClients.Grap(h.GetValue())
				if c == nil {
					continue
				}
				sck := c.GetSocket()
				slf.NetClients.Release(c)
				network.OperClose(sck)
			}

			for {
				time.Sleep(time.Duration(500) * time.Microsecond)
				if slf.NetClients.Size() <= 0 {
					break
				}

				slf.LogInfo("Service The remaining %d connections need to be closed", slf.NetClients.Size())
				ick++
				if ick > 6 {
					break
				}
			}
		}
	}
	slf.NetListen.Close()
	slf.LogInfo("Service Stoped")
}

//OnAccept Receive connection event
func (slf *NetListenService) OnAccept(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	accepter := message.(*network.NetAccept)
	if slf.NetClients.Size()+1 > slf.MaxClient {
		slf.LogWarning("OnAccept: client fulled:%d", slf.NetClients.Size())
		network.OperClose(accepter.Handle)
		return
	}

	c := slf.NetClients.Allocer().New()
	if c == nil {
		slf.LogError("OnAccept: client closed: insufficient memory")
		network.OperClose(accepter.Handle)
		return
	}

	c.SetSocket(accepter.Handle)
	c.SetAddr(accepter.Addr.String() + strconv.Itoa(accepter.Port))

	_, err := slf.NetClients.Occupy(c)
	if err != nil {
		slf.LogError("OnAccept: client closed: %v, %d-%s:%d",
			err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		slf.NetClients.Allocer().Delete(c)
		network.OperClose(accepter.Handle)
		return
	}

	network.OperOpen(accepter.Handle)
	network.OperSetKeep(accepter.Handle, slf.ClientKeep)

	if err = slf.NetDeleate.Handshake(c); err != nil {
		slf.LogError("OnAccept: client fail:%s", err)
	}

	c.GetStat().UpdateOnline(timer.Now())

	slf.NetClients.Release(c)

	slf.LogDebug("OnAccept: client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

//OnRecv Receiving data events
func (slf *NetListenService) OnRecv(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	defer slf.LogDebug("onRecv: complete")

	wrap := message.(*network.NetChunk)
	c := slf.NetClients.GrapSocket(wrap.Handle)
	if c == nil {
		slf.LogError("OnRecv: No target [%d] client object was found", wrap.Handle)
		return
	}

	defer slf.NetClients.Release(c)

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {
		space = c.GetRecvBuffer().Cap() - c.GetRecvBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = c.GetRecvBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				slf.LogError("OnRecv Write: error %+v socket %d", err, wrap.Handle)
				network.OperClose(wrap.Handle)
				break
			}

			pos += space
			writed += space

			c.GetStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = slf.NetDeleate.Analysis(context, slf, c)
			if err != nil {
				if err == ErrAnalysisSuccess {
					continue
				} else if err != ErrAnalysisProceed {
					slf.LogError("OnRecv: error %+v socket %d closing client", err, wrap.Handle)
					network.OperClose(wrap.Handle)
					return
				}
			}

			if writed >= len(wrap.Data) {
				return
			}

			break
		}
	}
}

//OnClose Close connection event
func (slf *NetListenService) OnClose(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	closer := message.(*network.NetClose)
	slf.LogDebug("close socket:%d", closer.Handle)
	c := slf.NetClients.GrapSocket(closer.Handle)
	if c == nil {
		slf.LogError("close unfind map-id socket %d", closer.Handle)
		return
	}

	defer slf.NetClients.Release(c)

	hClose := c.GetID()

	slf.NetClients.Erase(hClose)

	if err := slf.NetDeleate.UnOnlineNotification(hClose); err != nil {
		slf.LogDebug("closed client Notification %+v", err)
	}

	slf.LogDebug("closed client: %+v", hClose)
}

//Shutdown Termination of service
func (slf *NetListenService) Shutdown() {
	if slf.NetListen != nil {
		slf.NetListen.Close()
	}

	slf.Service.Shutdown()
}

//LogInfo Log information
func (slf *NetListenService) LogInfo(frmt string, args ...interface{}) {
	slf.Service.LogInfo(slf.getDesc()+frmt, args...)
}

//LogError Record error log information
func (slf *NetListenService) LogError(frmt string, args ...interface{}) {
	slf.Service.LogError(slf.getDesc()+frmt, args...)
}

//LogDebug Record debug log information
func (slf *NetListenService) LogDebug(frmt string, args ...interface{}) {
	slf.Service.LogDebug(slf.getDesc()+frmt, args...)
}

//LogTrace Record trace log information
func (slf *NetListenService) LogTrace(frmt string, args ...interface{}) {
	slf.Service.LogTrace(slf.getDesc()+frmt, args...)
}

//LogWarning Record warning log information
func (slf *NetListenService) LogWarning(frmt string, args ...interface{}) {
	slf.Service.LogWarning(slf.getDesc()+frmt, args...)
}
