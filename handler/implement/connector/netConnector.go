package connector

import (
	"errors"
	"time"

	libnet "github.com/yamakiller/magicLibs/net"
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

//Options doc
type Options struct {
	UID                int64
	Sock               net.INetConnection
	ReceiveBuffer      net.INetBuffer
	ReceiveDecoder     net.INetDecoder
	ReceiveOutChanSize int
	AsyncError         func(error)
	AsyncComplete      func(int32)
	AsyncClosed        func(int64)
}

var defaultOptions = Options{}

// Option is a function on the options for a connection.
type Option func(*Options) error

// Set uid
func SetUID(uid int64) Option {
	return func(o *Options) error {
		o.UID = uid
		return nil
	}
}

// Set Socket Handle
func SetSocket(s net.INetConnection) Option {
	return func(o *Options) error {
		o.Sock = s
		return nil
	}
}

func SetReceiveBuffer(b net.INetBuffer) Option {
	return func(o *Options) error {
		o.ReceiveBuffer = b
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

func SetAsyncClosed(f func(int64)) Option {
	return func(o *Options) error {
		o.AsyncClosed = f
		return nil
	}
}

func Spawn(options ...Option) (*NetConnector, error) {
	c := &NetConnector{_opts: defaultOptions}
	for _, opt := range options {
		if err := opt(&c._opts); err != nil {
			return nil, err
		}
	}

	if c._opts.Sock == nil {
		return nil, errors.New("NetConnector sock is null")
	}

	if c._opts.ReceiveBuffer == nil {
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
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
	slf.RegisterMethod(&network.NetChunk{}, slf.onRecv)
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
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
	slf._addr = addr
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
	slf._opts.Sock.Close()
	ick := 0
	for slf._status != UnConnected {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}
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

func (slf *NetConnector) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogDebug("Stoped: Socket-%d", slf._opts.Sock.GetSocket())
	slf._opts.Sock.WithSocket(libnet.INVALIDSOCKET)
	slf.Service.Stoped(context, sender, message)
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
	return
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

		space = slf._opts.ReceiveBuffer.Cap() - slf._opts.ReceiveBuffer.Len()
		wby = len(wrap.Data) - writed

		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = slf._opts.ReceiveBuffer.Write(wrap.Data[pos : pos+space])
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

//GetBufferCap doc
//@Summary Returns Recvice buffer cap
//@Method GetBufferCap
//@Return int
func (slf *NetConnector) GetBufferCap() int {
	return slf._opts.ReceiveBuffer.Cap()
}

//GetBufferLen doc
//@Summary Returns Recvice buffer data length
//@Method GetBufferLen
//@Return int
func (slf *NetConnector) GetBufferLen() int {
	return slf._opts.ReceiveBuffer.Len()
}

//GetBufferBytes doc
//@Summary Returns Recvice buffer all data
//@Method GetBufferBytes
//@Return []byte
func (slf *NetConnector) GetBufferBytes() []byte {
	return slf._opts.ReceiveBuffer.Bytes()
}

//ClearBuffer doc
//@Summary Clear Recvice buffer data
//@Method ClearBuffer
func (slf *NetConnector) ClearBuffer() {
	slf._opts.ReceiveBuffer.Clear()
}

//TrunBuffer doc
//@Summary Clear Recvice buffer n size data
//@Method TrunBuffer
//@Param int
func (slf *NetConnector) TrunBuffer(n int) {
	slf._opts.ReceiveBuffer.Trun(n)
}

//WriteBuffer doc
//@Summary Write Recvice buffer data
//@Method WriteBuffer
//@Return int
//@Return error
func (slf *NetConnector) WriteBuffer(b []byte) (int, error) {
	return slf._opts.ReceiveBuffer.Write(b)
}

//ReadBuffer doc
//@Summary Read Recvice buffer data
//@Method ReadBuffer
//@Param int read buffer size
//@Return []byte
func (slf *NetConnector) ReadBuffer(n int) []byte {
	return slf._opts.ReceiveBuffer.Read(n)
}

//SendTo doc
//@Summary Send data to network
//@Method SendTo
//@Param []byte send is data
//@Return error
func (slf *NetConnector) SendTo(data []byte) error {
	return network.OperWrite(slf._opts.Sock.GetSocket(), data, len(data))
}

//IsConnected doc
//@Summary is connected
//@Method IsConnected
//@Return bool
func (slf *NetConnector) IsConnected() bool {
	if slf._status == UnConnected {
		return false
	}
	return true
}

//OnClose doc
//@Summary  proccess closed connection events
//@Method OnClose
//@Param actor.Context  this actor context
//@Param *actor.PID     this message sender pid
//@Param interface{}    message
func (slf *NetConnector) OnClose(context actor.Context, sender *actor.PID, message interface{}) {
	slf._opts.ReceiveBuffer.Clear()
	slf._status = UnConnected
	if slf._opts.AsyncClosed != nil {
		slf._opts.AsyncClosed(slf._opts.UID)
	}
}
