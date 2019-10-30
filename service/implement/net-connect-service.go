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
func (slf *NetConnectService) Init() {
	slf.Service.Init()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&NetConnectEvent{}, slf.onConnection)
	slf.RegisterMethod(&network.NetChunk{}, slf.onRecv)
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
}

//Started Turn on network connect service
func (slf *NetConnectService) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf.Assignment(context)
	slf.LogInfo("Service Startup address:%s read-buffer-limit:%d chan-buffer-size:%d",
		slf.Target.GetAddr(),
		slf.Handle.GetRecvBufferLimit(),
		slf.Target.GetOutSize())
	slf.Service.Started(context, sender, message)
	slf.LogInfo("Service Startup completed")
}

//Stopping Out of service
func (slf *NetConnectService) Stopping(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogInfo("[%s] %s Connection Service Stoping %s",
		slf.Handle.Name(),
		slf.Name(),
		slf.Target.GetAddr())
	slf.isShutdown = false
	slf.Handle.Close()
	slf.NetMethod.Clear()
	slf.LogInfo("Connection Service Stoped %s", slf.Target.GetAddr())
}

//IsShutdown Whether the service has been terminated
func (slf *NetConnectService) IsShutdown() bool {
	return slf.isShutdown
}

//AutoConnect  auto connect
func (slf *NetConnectService) AutoConnect(context actor.Context) error {
	err := slf.Handle.Connection(context, slf.Target.GetAddr(), slf.Target.GetOutSize())
	if err != nil {
		goto unend
	}

	err = slf.Deleate.Connected(context, slf)
	if err != nil {
		slf.Handle.Close()
		goto unend
	}
	return nil
unend:
	slf.Target.SetEtat(UnConnected)
	return err
}

//onConnection Request connection event
func (slf *NetConnectService) onConnection(context actor.Context, sender *actor.PID, message interface{}) {
	//t := message.(*NetConnectEvent)
	slf.LogInfo("onConnection: %s", slf.Target.GetAddr())
	err := slf.AutoConnect(context)
	if err != nil {
		slf.LogError("onConnection: fail-%+v", err)
	}
}

//OnRecv Connection read data
func (slf *NetConnectService) onRecv(context actor.Context, sender *actor.PID, message interface{}) {
	defer slf.LogDebug("onRecv complete")
	wrap := message.(*network.NetChunk)
	if wrap.Handle != slf.Handle.Socket() {
		slf.LogDebug("[%d:%d]Discard the data because this data is the current connection authorization data.",
			wrap.Handle,
			slf.Handle.Socket())
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
		if slf.isShutdown {
			break
		}

		space = slf.Handle.GetRecvBufferLimit() - slf.Handle.GetRecvBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = slf.Handle.GetRecvBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				slf.Handle.Close()
				break
			}

			pos += space
			writed += space

			slf.Handle.GetDataStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = slf.Deleate.Analysis(context, slf)
			if err != nil {
				if err == ErrAnalysisSuccess {
					continue
				} else if err != ErrAnalysisProceed {
					slf.Handle.Close()
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
func (slf *NetConnectService) OnClose(context actor.Context, sender *actor.PID, message interface{}) {
	//Release buffer resources
	slf.Handle.GetRecvBuffer().Reset()
	slf.Target.SetEtat(UnConnected)
}

// Shutdown : Proactively shut down the service
func (slf *NetConnectService) Shutdown() {
	slf.isShutdown = true
	if slf.Handle.Socket() != 0 {
		network.OperClose(slf.Handle.Socket())
	}
	slf.Service.Shutdown()
}

func (slf *NetConnectService) getDesc() string {
	return fmt.Sprintf("[%s] %s ", slf.Handle.Name(), slf.Name())
}

//LogInfo Log information
func (slf *NetConnectService) LogInfo(frmt string, args ...interface{}) {
	slf.Service.LogInfo(slf.getDesc()+frmt, args...)
}

//LogError Record error log information
func (slf *NetConnectService) LogError(frmt string, args ...interface{}) {
	slf.Service.LogError(slf.getDesc()+frmt, args...)
}

//LogDebug Record debug log information
func (slf *NetConnectService) LogDebug(frmt string, args ...interface{}) {
	slf.Service.LogDebug(slf.getDesc()+frmt, args...)
}

//LogTrace Record trace log information
func (slf *NetConnectService) LogTrace(frmt string, args ...interface{}) {
	slf.Service.LogTrace(slf.getDesc()+frmt, args...)
}

//LogWarning Record warning log information
func (slf *NetConnectService) LogWarning(frmt string, args ...interface{}) {
	slf.Service.LogWarning(slf.getDesc()+frmt, args...)
}
