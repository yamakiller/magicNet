package implement

import (
	"errors"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
	"github.com/yamakiller/magicNet/handler/net"
	"github.com/yamakiller/magicNet/network"
)

var (
	//
	Idel = 0
	//UnConnected Not connected or failed to connect
	UnConnected = 1
	//Connecting Connecting target
	Connecting = 2
	//Verify Performing a certification login
	Verify = 3
	//Connected Connection has been completed
	Connected = 4
)

var (
	ErrNetConnecting = errors.New("network is connecting")
)

type netConnectEvent struct {
	_addr          string
	_asyncError    func(err error)
	_asyncComplete func(sock int32)
}

type NetConnector struct {
	handler.Service

	_sock    net.INetConnection
	_status  int
	_outSize int
}

func (slf *NetConnector) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&netConnectEvent{}, slf.onConnection)
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
}

func (slf *NetConnector) Connection(addr string, timeout int, asyncError func(err error), asyncComplete func(sock int32)) error {
	ick := 0
	for slf._status == Idel {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	if slf._status != UnConnected {
		return ErrNetConnecting
	}

	slf._status = Connecting
	actor.DefaultSchedulerContext.Send(slf.GetPID(), &netConnectEvent{addr, asyncError, asyncComplete})
	return nil
}

func (slf *NetConnector) Shutdown() {

	ick := 0
	for slf._status != UnConnected {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	if slf._sock.GetSocket() != 0 {
		network.OperClose(slf._sock.GetSocket())
	}

	slf.Service.Shutdown()
}

func (slf *NetConnector) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf.Service.Started(context, sender, message)
	slf._status = UnConnected
}

func (slf *NetConnector) Stopping(context actor.Context, sender *actor.PID, message interface{}) {

	slf.Started(context, sender, message)
}

func (slf *NetConnector) onConnection(context actor.Context, sender *actor.PID, message interface{}) {
	req := message.(*netConnectEvent)
	err := slf._sock.Connection(context, req._addr, slf._outSize)
	if err != nil {
		goto unend
	}

unend:
	slf._status = UnConnected
	if err != nil && req._asyncError != nil {
		req._asyncError(err)
	}
	return
}

//NetConnectStatus Connection Status
/*type NetConnectStatus int32

var (
	//UnConnected Not connected or failed to connect
	UnConnected = NetConnectStatus(0)
	//Connecting Connecting target
	Connecting = NetConnectStatus(1)
	//Verify Performing a certification login
	Verify = NetConnectStatus(2)
	//Connected Connection has been completed
	Connected = NetConnectStatus(3)
)

//INetConnectTarget Connection target interface
type INetConnectTarget interface {
	GetAddr() string
	GetOutSize() int
	SetStatus(stat NetConnectStatus)
	GetStatus() NetConnectStatus
}

//NetConnectEvent Connection event
type NetConnectEvent struct {
	Target INetConnectTarget
}

//INetConnectDeleate Commission
type INetConnectDeleate interface {
	Connected(context actor.Context, nets *NetConnectService) error
	Decode(context actor.Context, nets *NetConnectService) error
}

//NetConnectService Internet connection service
type NetConnectService struct {
	handler.Service
	Handle     net.INetConnection
	Deleate    INetConnectDeleate
	Target     INetConnectTarget
	NetMethod  NetMethodDispatch
	isShutdown bool
}

//Initial Initialize the network listening service
func (slf *NetConnectService) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&NetConnectEvent{}, slf.onConnection)
	slf.RegisterMethod(&network.NetChunk{}, slf.onRecv)
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
}

//Started Turn on network connect service
func (slf *NetConnectService) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogInfo("Service Startup address:%s read-buffer-limit:%d chan-buffer-size:%d",
		slf.Target.GetAddr(),
		slf.Handle.GetReceiveBufferLimit(),
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
	slf.Target.SetStatus(UnConnected)
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

		space = slf.Handle.GetReceiveBufferLimit() - slf.Handle.GetReceiveBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = slf.Handle.GetReceiveBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				slf.Handle.Close()
				break
			}

			pos += space
			writed += space

			slf.Handle.GetStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = slf.Deleate.Decode(context, slf)
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
	slf.Handle.GetReceiveBuffer().Reset()
	slf.Target.SetStatus(UnConnected)
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
}*/
