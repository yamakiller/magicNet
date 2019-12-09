package connector

import (
	"errors"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
	"github.com/yamakiller/magicNet/handler/net"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/timer"
)

const (
	//
	Idle = 0
	//UnConnected Not connected or failed to connect
	UnConnected = 1
	//Connecting Connecting target
	Connecting = 2
	//Connected Connection has been completed
	Connected = 3
)

var (
	//ErrNetConnecting Being connected error
	ErrNetConnecting = errors.New("network is connecting")
)

type netConnectEvent struct {
}

type Options struct {
	Sock               net.INetConnection
	Receive            net.INetBuffer
	ReceiveDecoder     net.INetDecoder
	ReceiveOutChanSize int
	AsyncError         func(error)
	AsyncComplete      func(int32)
}

var DefaultOptions = Options{}

// Option is a function on the options for a connection.
type Option func(*Options) error

func SetSocket(s net.INetConnection) Option {
	return func(o *Options) error {
		o.Sock = s
		return nil
	}
}

func SetReceive(b net.INetBuffer) Option {
	return func(o *Options) error {
		o.Receive = b
		return nil
	}
}

func SetReceiveDecoder(d net.INetDecoder) Option {
	return func(o *Options) error {
		o.ReceiveDecoder = d
		return nil
	}
}

func SetReceiveOutChanSize(ocs int) Option {
	return func(o *Options) error {
		o.ReceiveOutChanSize = ocs
		return nil
	}
}

func SetAsyncError(f func(error)) Option {
	return func(o *Options) error {
		o.AsyncError = f
		return nil
	}
}

func SetAsyncComplete(f func(int32)) Option {
	return func(o *Options) error {
		o.AsyncComplete = f
		return nil
	}
}

func Spawn(options ...Option) (*NetConnector, error) {
	c := &NetConnector{_opts: DefaultOptions}
	for _, opt := range options {
		if err := opt(&c._opts); err != nil {
			return nil, err
		}
	}

	if c._opts.Sock == nil {
		return nil, errors.New("NetConnector sock is null")
	}

	if c._opts.Receive == nil {
		return nil, errors.New("NetConnector receive buffer is null")
	}

	if c._opts.ReceiveDecoder == nil {
		return nil, errors.New("NetConnector receive decoder is null")
	}

	if c._opts.ReceiveOutChanSize <= 0 {
		return nil, errors.New("NetConnector receive out chan size is 0")
	}

	return c, nil
}

//NetConnector doc
//@Summary network connector
//@Struct NetConnector
//@Inherit Service
//@Member  net.INetConnection  connection interface
//@Member  net.INetBuffer      connection receive buffer
//@Member  net.INetDecoder     connection receive decoder
//@Member  int                 connection status
//@Member  int                 connection receive out chan size
//@Member  int64               connection receive bytes
//@Member  int64               connection receive last time
//@Member  int64               connection sendout bytes
//@Member  int64               connection sendout last time
type NetConnector struct {
	handler.Service

	_addr            string
	_opts            Options
	_status          int
	_receiveBytes    int64
	_receiveLastTime int64
	_sendoutBytes    int64
	_sendoutLastTime int64
}

//Initial doc
//@Summary initialization network connector
//@Method Initial
func (slf *NetConnector) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&netConnectEvent{}, slf.onConnection)
	slf.RegisterMethod(&network.NetChunk{}, slf.onRecv)
	slf.RegisterMethod(&actor.Started{}, slf.Started)
}

//Connection doc
//@Summary connection
//@Method Connection
//@Param  string   			connection address [ip:port]
//@Param  int				connection Receive chan size
//@Param  func(err error)   connection async error function
//@Param  func(sock int32)  connection async complete function
//@Return error
func (slf *NetConnector) Connection(addr string) error {
	ick := 0
	for slf._status == Idle {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	if slf._status != UnConnected {
		return ErrNetConnecting
	}

	now := int64(timer.Now())
	slf._receiveBytes = 0
	slf._sendoutBytes = 0
	slf._receiveLastTime = now
	slf._sendoutLastTime = now
	slf._status = Connecting
	actor.DefaultSchedulerContext.Send(slf.GetPID(), &netConnectEvent{})
	return nil
}

//Shutdown doc
//@Summary shutdown connector
//@Method Shutdown
func (slf *NetConnector) Shutdown() {

	ick := 0
	for slf._status != UnConnected {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	slf._opts.Sock.Close()
	slf.Service.Shutdown()
}

//Started doc
//@Summary Started event
//@Method Started
//@Param  actor.Context   context
//@Param  *actor.PID      sender
//@Param  interface{}     message
func (slf *NetConnector) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf.Service.Started(context, sender, message)
	slf._status = UnConnected
}

func (slf *NetConnector) onConnection(context actor.Context, sender *actor.PID, message interface{}) {
	err := slf._opts.Sock.Connection(context, slf._addr, slf._opts.ReceiveOutChanSize)
	if err != nil {
		goto unend
	}

	slf._status = Connected
	if slf._opts.AsyncComplete != nil {
		slf._opts.AsyncComplete(slf._opts.Sock.GetSocket())
	}
unend:
	slf._status = UnConnected
	if err != nil && slf._opts.AsyncError != nil {
		slf._opts.AsyncError(err)
	}
	return
}

func (slf *NetConnector) onRecv(context actor.Context, sender *actor.PID, message interface{}) {
	wrap := message.(*network.NetChunk)
	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {

		space = slf._opts.Receive.Cap() - slf._opts.Receive.Len()
		wby = len(wrap.Data) - writed

		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = slf._opts.Receive.Write(wrap.Data[pos : pos+space])
			if err != nil {
				slf._opts.Sock.Close()
				break
			}

			pos += space
			writed += space

			slf._receiveBytes += int64(space)
			slf._receiveLastTime = int64(timer.Now())
		}

		for {
			// Decomposition of Packets
			err = slf._opts.ReceiveDecoder(context, slf)
			if err != nil {
				if err == net.ErrAnalysisSuccess {
					continue
				} else if err != net.ErrAnalysisProceed {
					slf._opts.Sock.Close()
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

//OnClose doc
//@Summary  proccess closed connection events
//@Method OnClose
//@Param actor.Context  this actor context
//@Param *actor.PID     this message sender pid
//@Param interface{}    message
func (slf *NetConnector) OnClose(context actor.Context, sender *actor.PID, message interface{}) {
	slf._opts.Receive.Clear()
	slf._status = UnConnected
}
