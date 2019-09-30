package implement

import (
	"fmt"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/net"
	"github.com/yamakiller/magicNet/timer"
)

//NetConnectEtat Connection Status
type NetConnectEtat int32

var (
	//UnConnected Not connected or failed to connect
	UnConnected = NetConnectEtat(0)
	//Connecting Connecting target
	Connecting = NetConnectEtat(1)
	//Verify Performing a certification login
	Verify = NetConnectEtat(2)
	//Connected Connection has been completed
	Connected = NetConnectEtat(3)
)

//INetConnectTarget Connection target interface
type INetConnectTarget interface {
	GetName() string
	GetAddr() string
	GetOutSize() int
	IsTimeout() uint64
	RestTimeout()
	SetEtat(stat NetConnectEtat)
	GetEtat() NetConnectEtat
}

//NetConnectEvent Connection event
type NetConnectEvent struct {
	Target INetConnectTarget
}

//INetConnectDeleate Commission
type INetConnectDeleate interface {
	Connected(context actor.Context, nets *NetConnectService) error
	Analysis(context actor.Context, nets *NetConnectService) error
}

//NetConnectService Internet connection service
type NetConnectService struct {
	service.Service
	Handle     net.INetConnection
	Deleate    INetConnectDeleate
	Target     INetConnectTarget
	NetMethod  NetMethodDispatch
	isShutdown bool
}

//Init Initialize the network listening service
func (nets *NetConnectService) Init() {
	nets.Service.Init()
	nets.RegisterMethod(&actor.Started{}, nets.Started)
	nets.RegisterMethod(&actor.Stopping{}, nets.Stopping)
	nets.RegisterMethod(&NetConnectEvent{}, nets.onConnection)
	nets.RegisterMethod(&network.NetChunk{}, nets.onRecv)
	nets.RegisterMethod(&network.NetClose{}, nets.OnClose)
}

//Started Turn on network connect service
func (nets *NetConnectService) Started(context actor.Context, message interface{}) {
	nets.Assignment(context)
	nets.LogInfo("Service Startup address:%s read-buffer-limit:%d chan-buffer-size:%d",
		nets.Target.GetAddr(),
		nets.Handle.GetRecvBufferLimit(),
		nets.Target.GetOutSize())
	nets.Service.Started(context, message)
	nets.LogInfo("Service Startup completed")
}

//Stopping Out of service
func (nets *NetConnectService) Stopping(context actor.Context, message interface{}) {
	nets.LogInfo("[%s] %s Connection Service Stoping %s",
		nets.Handle.Name(),
		nets.Name(),
		nets.Target.GetAddr())
	nets.isShutdown = false
	nets.Handle.Close()
	nets.NetMethod.Clear()
	nets.LogInfo("Connection Service Stoped %s", nets.Target.GetAddr())
}

//IsShutdown Whether the service has been terminated
func (nets *NetConnectService) IsShutdown() bool {
	return nets.isShutdown
}

//AutoConnect  auto connect
func (nets *NetConnectService) AutoConnect(context actor.Context) error {
	err := nets.Handle.Connection(context, nets.Target.GetAddr(), nets.Target.GetOutSize())
	if err != nil {
		goto unend
	}

	err = nets.Deleate.Connected(context, nets)
	if err != nil {
		nets.Handle.Close()
		goto unend
	}
	return nil
unend:
	nets.Target.SetEtat(UnConnected)
	return err
}

//onConnection Request connection event
func (nets *NetConnectService) onConnection(context actor.Context, message interface{}) {
	//t := message.(*NetConnectEvent)
	nets.LogInfo("onConnection: %s", nets.Target.GetAddr())
	err := nets.AutoConnect(context)
	if err != nil {
		nets.LogError("onConnection: fail-%+v", err)
	}
}

//OnRecv Connection read data
func (nets *NetConnectService) onRecv(context actor.Context, message interface{}) {
	defer nets.LogDebug("onRecv complete")
	wrap := message.(*network.NetChunk)
	if wrap.Handle != nets.Handle.Socket() {
		nets.LogDebug("[%d:%d]Discard the data because this data is the current connection authorization data.",
			wrap.Handle,
			nets.Handle.Socket())
		return
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {
		if nets.isShutdown {
			break
		}

		space = nets.Handle.GetRecvBufferLimit() - nets.Handle.GetRecvBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = nets.Handle.GetRecvBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				nets.Handle.Close()
				break
			}

			pos += space
			writed += space

			nets.Handle.GetDataStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = nets.Deleate.Analysis(context, nets)
			if err != nil {
				if err == ErrAnalysisSuccess {
					continue
				} else if err != ErrAnalysisProceed {
					nets.Handle.Close()
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

//OnClose Handling closed connection events
func (nets *NetConnectService) OnClose(context actor.Context, message interface{}) {
	//Release buffer resources
	nets.Handle.GetRecvBuffer().Reset()
	nets.Target.SetEtat(UnConnected)
}

// Shutdown : Proactively shut down the service
func (nets *NetConnectService) Shutdown() {
	nets.isShutdown = true
	if nets.Handle.Socket() != 0 {
		network.OperClose(nets.Handle.Socket())
	}
	nets.Service.Shutdown()
}

func (nets *NetConnectService) getDesc() string {
	return fmt.Sprintf("[%s] %s ", nets.Handle.Name(), nets.Name())
}

//LogInfo Log information
func (nets *NetConnectService) LogInfo(frmt string, args ...interface{}) {
	nets.Service.LogInfo(nets.getDesc()+frmt, args...)
}

//LogError Record error log information
func (nets *NetConnectService) LogError(frmt string, args ...interface{}) {
	nets.Service.LogError(nets.getDesc()+frmt, args...)
}

//LogDebug Record debug log information
func (nets *NetConnectService) LogDebug(frmt string, args ...interface{}) {
	nets.Service.LogDebug(nets.getDesc()+frmt, args...)
}

//LogTrace Record trace log information
func (nets *NetConnectService) LogTrace(frmt string, args ...interface{}) {
	nets.Service.LogTrace(nets.getDesc()+frmt, args...)
}

//LogWarning Record warning log information
func (nets *NetConnectService) LogWarning(frmt string, args ...interface{}) {
	nets.Service.LogWarning(nets.getDesc()+frmt, args...)
}
