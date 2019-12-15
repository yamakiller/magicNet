package listener

import (
	"errors"
	"strconv"
	"time"

	"github.com/yamakiller/magicNet/timer"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
	"github.com/yamakiller/magicNet/handler/net"
	"github.com/yamakiller/magicNet/network"
)

const (
	//Idle listener status Idle
	Idle = 0
	//UnListen listener status unlisten
	UnListen = 1
	//Listening listener status listening
	Listening = 2
	//Listened listener status listened
	Listened = 3
)

var (
	//ErrNetListened Being listened error
	ErrNetListened = errors.New("network is listened")
)

type AsyncErrorFunc func(error)
type AsyncCompleteFunc func(int32)
type AsyncAcceptFunc func(net.INetClient) error
type AsyncClosedFunc func(uint64) error

type netListenEvent struct {
}

//Options doc
type Options struct {
	Sock           net.INetListener
	CSGroup        net.INetClientGroup
	KeepTime       int
	OutCChanSize   int
	ReceiveDecoder net.INetDecoder

	AsyncError    AsyncErrorFunc
	AsyncComplete AsyncCompleteFunc
	AsyncAccept   AsyncAcceptFunc
	AsyncClosed   AsyncClosedFunc
}

//DefaultOptions doc
var DefaultOptions = Options{}

//Option is a function on the options for a listen.
type Option func(*Options) error

//SetListener doc
//@Summary Set the listener handle object
//@Method
//@Param net.INetListener Listening handle/TCP/UDP/KCP/WebSocket
//@Return Option
func SetListener(s net.INetListener) Option {
	return func(o *Options) error {
		o.Sock = s
		return nil
	}
}

//SetClientGroups doc
//@Summary Set up a connection client management group
//@Method  SetClientGroups
//@Param   net.INetClientGroup Management Group Object
//@Return  Option
func SetClientGroups(b net.INetClientGroup) Option {
	return func(o *Options) error {
		o.CSGroup = b
		return nil
	}
}

//SetClientOutChanSize doc
//@Summary Set the connection client transaction pipeline buffer size
//@Method SetClientOutChanSize
//@Param  int Pipe buffer size
func SetClientOutChanSize(ch int) Option {
	return func(o *Options) error {
		o.OutCChanSize = ch
		return nil
	}
}

//SetClientDecoder doc
//@Summary Set the connection client data decoder
//@Method SetClientDecoder
//@Param  net.INetDecoder decoder
//@Return Option
func SetClientDecoder(d net.INetDecoder) Option {
	return func(o *Options) error {
		o.ReceiveDecoder = d
		return nil
	}
}

//SetClientKeepTime doc
//@Summary Set the heartbeat interval of the connected client in milliseconds
//@Param   int Interval time in milliseconds
//@Return  Option
func SetClientKeepTime(tm int) Option {
	return func(o *Options) error {
		o.KeepTime = tm
		return nil
	}
}

//SetAsyncError doc
//@Summary Set the callback function to listen for asynchronous errors
//@Method SetAsyncError
//@Param  func(error) Callback
//@Return Option
func SetAsyncError(f AsyncErrorFunc) Option {
	return func(o *Options) error {
		o.AsyncError = f
		return nil
	}
}

//SetAsyncComplete doc
//@Summary Set the callback completion asynchronous callback function
//@Method SetAsyncComplete
//@Param  func(int32) Callback
//@Return Option
func SetAsyncComplete(f AsyncCompleteFunc) Option {
	return func(o *Options) error {
		o.AsyncComplete = f
		return nil
	}
}

//SetAsyncAccept doc
//@Summary  Set listen accept asynchronous callback function
//@Method   SetAsyncAccept
//@Param    func(net.INetClient) error  Callback
//@Return   Option
func SetAsyncAccept(f AsyncAcceptFunc) Option {
	return func(o *Options) error {
		o.AsyncAccept = f
		return nil
	}
}

//SetAsyncClose doc
//@Summary Set the client to close the asynchronous callback function
//@Method Close
//@Param  func(uint64) error Callback
//@Return Option
func SetAsyncClosed(f AsyncClosedFunc) Option {
	return func(o *Options) error {
		o.AsyncClosed = f
		return nil
	}
}

//Spawn doc
//@Summary Create a listening service object
//@Method Spawn
//@Param  ...Option Setting parameters
//@Return *NetListener Listening service object
//@Return error
func Spawn(options ...Option) (*NetListener, error) {
	c := &NetListener{_opts: DefaultOptions}
	for _, opt := range options {
		if err := opt(&c._opts); err != nil {
			return nil, err
		}
	}

	if c._opts.Sock == nil {
		return nil, errors.New("NetListener sock is null")
	}

	if c._opts.CSGroup == nil {
		return nil, errors.New("NetListener client group is null")
	}

	if c._opts.ReceiveDecoder == nil {
		return nil, errors.New("NetListener receive decoder is null")
	}

	if c._opts.OutCChanSize <= 0 {
		return nil, errors.New("NetListener client receive out chan is 0")
	}

	return c, nil
}

//NetListener doc
//@Summary network listen
//@Struct NetListener
//@Inherit Service
type NetListener struct {
	handler.Service

	_addr   string
	_opts   Options
	_status int
}

//Initial Initialize the network listening service
func (slf *NetListener) Initial() {
	slf.Service.Initial()
	slf._opts.CSGroup.Initial()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&netListenEvent{}, slf.onListen)
	slf.RegisterMethod(&network.NetAccept{}, slf.onAccept)
	slf.RegisterMethod(&network.NetChunk{}, slf.onRecv)
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
}

//Listen doc
//@Summary Listener
//@Method Listen
//@Param  string   			listen address [ip:port]
//@Param  int               listen client keepTime millsecond
//@Param  int				listen client Receive chan size
//@Param  func(err error)   connection async error function
//@Param  func(sock int32)  connection async complete function
//@Return error
func (slf *NetListener) Listen(addr string) error {
	ick := 0
	for slf._status == Idle {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	if slf._status != UnListen {
		return ErrNetListened
	}

	slf._addr = addr
	slf._status = Listening
	actor.DefaultSchedulerContext.Send(slf.GetPID(), &netListenEvent{})

	return nil
}

//Shutdown doc
//@Summary Termination of service
//@Method Shutdown
func (slf *NetListener) Shutdown() {
	hs := slf._opts.CSGroup.GetHandles()
	if hs != nil && len(hs) > 0 {
		for slf._opts.CSGroup.Size() > 0 {
			ick := 0
			for i := 0; i < len(hs); i++ {
				c := slf._opts.CSGroup.Grap(hs[i])
				if c == nil {
					continue
				}

				sock := c.GetSocket()
				slf._opts.CSGroup.Release(c)
				network.OperClose(sock)
			}

			for {
				time.Sleep(time.Duration(500) * time.Microsecond)
				if slf._opts.CSGroup.Size() <= 0 {
					break
				}

				slf.LogDebug("Service The remaining %d connections need to be closed", slf._opts.CSGroup.Size())
				ick++
				if ick > 6 {
					break
				}
			}
		}
	}

	slf._opts.Sock.Close()
	slf.Service.Shutdown()
}

//Grap doc
//@Summary Grap Client
//@Method Grap
//@Param  uint64 client handle
//@Return net.INetClient
func (slf *NetListener) Grap(handle uint64) net.INetClient {
	return slf._opts.CSGroup.Grap(handle)
}

//GetClients doc
//@Summary Return all client handle
//@Method GetClients
//@Return []uint64
func (slf *NetListener) GetClients() []uint64 {
	return slf._opts.CSGroup.GetHandles()
}

//Release doc
//@Summary Return Grap Get Client
//@Param net.INetCLient
func (slf *NetListener) Release(c net.INetClient) {
	slf._opts.CSGroup.Release(c)
}

//Started doc
//@Summary Started event
//@Method Started
//@Param  actor.Context
//@Param *actor.PID
//@Param message
func (slf *NetListener) Started(context actor.Context, sender *actor.PID, message interface{}) {
	slf._status = UnListen
	slf.Service.Started(context, sender, message)
}

func (slf *NetListener) onListen(context actor.Context, sender *actor.PID, message interface{}) {
	err := slf._opts.Sock.Listen(context, slf._addr, slf._opts.OutCChanSize)
	if err != nil {
		slf._status = UnListen
		if slf._opts.AsyncError != nil {
			slf._opts.AsyncError(err)
		}
		return
	}

	slf._status = Listened
	if slf._opts.AsyncComplete != nil {
		slf._opts.AsyncComplete(slf._opts.Sock.GetSocket())
	}
}

//onAccept doc
//@Summary accept connection event
//@Method onAccept
//@Param  actor.Context
//@Param *actor.PID
//@Param message
func (slf *NetListener) onAccept(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	accepter := message.(*network.NetAccept)
	if slf._opts.CSGroup.Size()+1 > slf._opts.CSGroup.Cap() {
		slf.LogWarning("OnAccept: client fulled-%d", slf._opts.CSGroup.Size())
		network.OperClose(accepter.Handle)
		return
	}

	c := slf._opts.CSGroup.Allocer().New()
	if c == nil {
		slf.LogError("OnAccept: client closed-insufficient memory")
		network.OperClose(accepter.Handle)
		return
	}

	c.WithSocket(accepter.Handle)
	c.WithAddr(accepter.Addr.String() + strconv.Itoa(accepter.Port))

	id, err := slf._opts.CSGroup.Occupy(c)
	if err != nil {
		slf.LogError("OnAccept: client closed-%v, %d-%s:%d",
			err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		slf._opts.CSGroup.Allocer().Delete(c)
		network.OperClose(accepter.Handle)
		return
	}

	defer slf._opts.CSGroup.Release(c)
	if err = network.OperOpen(accepter.Handle); err != nil {
		slf._opts.CSGroup.Erase(id)
		slf.LogError("OnAccept: client open fail-%+v", err)
		return
	}

	if err = network.OperSetKeep(accepter.Handle, uint64(slf._opts.KeepTime)); err != nil {
		slf._opts.CSGroup.Erase(id)
		slf.LogError("OnAccept: client setkeep fail-%+v", err)
		return
	}

	c.UpdateOnline(int64(timer.Now()))
	if err = slf._opts.AsyncAccept(c); err != nil {
		slf._opts.CSGroup.Erase(id)
		slf.LogError("OnAccept: client fail-%+v", err)
		return
	}

	slf.LogDebug("OnAccept: client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

//onRecv doc
//@Summary onRecv connection event
//@Method Receive
//@Param  actor.Context
//@Param *actor.PID
//@Param message
func (slf *NetListener) onRecv(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	wrap := message.(*network.NetChunk)
	c := slf._opts.CSGroup.GrapSocket(wrap.Handle)
	if c == nil {
		slf.LogError("OnRecv: No target [%d] client object was found", wrap.Handle)
		return
	}

	defer slf._opts.CSGroup.Release(c)

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {
		space = c.GetBufferCap() - c.GetBufferLen()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = c.WriteBuffer(wrap.Data[pos : pos+space])
			if err != nil {
				slf.LogError("OnRecv Write: error %+v socket %d", err, wrap.Handle)
				network.OperClose(wrap.Handle)
				break
			}

			pos += space
			writed += space

			c.UpdateReceive(int64(timer.Now()), int64(space))
		}

		for {
			// Decomposition of Packets
			err = slf._opts.ReceiveDecoder(context, slf, c)
			if err != nil {
				if err == net.ErrAnalysisSuccess {
					continue
				} else if err != net.ErrAnalysisProceed {
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
func (slf *NetListener) OnClose(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	closer := message.(*network.NetClose)
	slf.LogDebug("close socket:%d", closer.Handle)
	c := slf._opts.CSGroup.GrapSocket(closer.Handle)
	if c == nil {
		slf.LogError("close unfind map-id socket %d", closer.Handle)
		return
	}

	defer slf._opts.CSGroup.Release(c)

	sockHandle := c.GetID()

	slf._opts.CSGroup.Erase(sockHandle)

	if slf._opts.AsyncClosed != nil {
		if err := slf._opts.AsyncClosed(sockHandle); err != nil {
			slf.LogError("closed client Notification %+v", err)
		}
	}
	slf.LogDebug("closed client: %+v", sockHandle)
}
