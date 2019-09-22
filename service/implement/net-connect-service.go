package implement

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/net"
	"github.com/yamakiller/magicNet/timer"
)

type NetConnStat int32

var (
	UnConnected = NetConnStat(0)
	Connecting  = NetConnStat(1)
	Connected   = NetConnStat(2)
)

type INetConnectionTarget interface {
	GetAddr() string
	GetOutSize() int
	SetStatus(stat NetConnStat)
}

type NetConnectionEvent struct {
	Target INetConnectionTarget
}

type NetConnForwadEvent struct {
	Wrap []byte
}

type INetConnectDeleate interface {
	Connected(context actor.Context, nets *NetConnectService) error
	Forawd(context actor.Context, nets *NetConnectService, message interface{}) error
	Analysis(context actor.Context, nets *NetConnectService) error
}

type NetConnectService struct {
	service.Service
	Handle     net.INetConnection
	Deleate    INetConnectDeleate
	Target     INetConnectionTarget
	isShutdown bool
}

//Init Initialize the network listening service
func (nets *NetConnectService) Init() {
	nets.Service.Init()
	nets.RegisterMethod(&actor.Started{}, nets.Started)
	nets.RegisterMethod(&actor.Stopped{}, nets.Stoped)
	nets.RegisterMethod(&NetConnectionEvent{}, nets.OnConnection)
	nets.RegisterMethod(&NetConnForwadEvent{}, nets.OnForwad)
	nets.RegisterMethod(&network.NetChunk{}, nets.OnRecv)
	nets.RegisterMethod(&network.NetClose{}, nets.OnClose)
}

//OnConnection Request connection event
func (nets *NetConnectService) OnConnection(context actor.Context, message interface{}) {
	t := message.(*NetConnectionEvent)
	err := nets.Handle.Connection(context, t.Target.GetAddr(), t.Target.GetOutSize())
	if err != nil {
		goto unend
	}

	err = nets.Deleate.Connected(context, nets)
	if err != nil {
		goto unend
	}
	return
unend:
	t.Target.SetStatus(UnConnected)
}

//OnRecv Connection read data
func (nets *NetConnectService) OnRecv(context actor.Context, message interface{}) {
	defer nets.LogDebug("onRecv complete")

	wrap := message.(*network.NetChunk)

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {
		space = nets.Handle.GetRecvBufferLimit() - nets.Handle.GetRecvBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = nets.Handle.GetRecvBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				network.OperClose(wrap.Handle)
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

//OnForwad send data
func (nets *NetConnectService) OnForwad(context actor.Context, message interface{}) {
	m, ok := message.(*NetConnForwadEvent)
	if !ok {
		return
	}

	err := nets.Deleate.Forawd(context, nets, m)
	if err != nil {
		return
	}

	//logger
}

//OnClose Handling closed connection events
func (nets *NetConnectService) OnClose(context actor.Context, message interface{}) {
	nets.Handle.Close()
	nets.Target.SetStatus(UnConnected)
}

// Shutdown : Proactively shut down the service
func (nets *NetConnectService) Shutdown() {
	if nets.Handle.Socket() != 0 {
		network.OperClose(nets.Handle.Socket())
	}
	nets.Service.Shutdown()
}
