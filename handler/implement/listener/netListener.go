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

type netListenEvent struct {
}

type Options struct {
	Sock           net.INetListener
	CSGroup        net.INetClientGroup
	KeepTime       int
	OutCChanSize   int
	ReceiveDecoder net.INetDecoder

	AsyncError    func(error)
	AsyncComplete func(int32)
	AsyncAccept   func(net.INetClient) error
	AsyncClose    func(uint64) error
}

var DefaultOptions = Options{}

// Option is a function on the options for a listen.
type Option func(*Options) error

func SetListener(s net.INetListener) Option {
	return func(o *Options) error {
		o.Sock = s
		return nil
	}
}

func SetClientGroups(b net.INetClientGroup) Option {
	return func(o *Options) error {
		o.CSGroup = b
		return nil
	}
}

func SetClientOutChanSize(ch int) Option {
	return func(o *Options) error {
		o.OutCChanSize = ch
		return nil
	}
}

func SetClientDecoder(d net.INetDecoder) Option {
	return func(o *Options) error {
		o.ReceiveDecoder = d
		return nil
	}
}

func SetClientKeepTime(tm int) Option {
	return func(o *Options) error {
		o.KeepTime = tm
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

func SetAsyncAccept(f func(net.INetClient) error) Option {
	return func(o *Options) error {
		o.AsyncAccept = f
		return nil
	}
}

func SetAsyncClose(f func(uint64) error) Option {
	return func(o *Options) error {
		o.AsyncClose = f
		return nil
	}
}

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
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
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
	for slf._status == UnListen {
		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(2) * time.Millisecond)
		}
	}

	if slf._status != UnListen {
		return ErrNetListened
	}

	slf._status = Listening
	actor.DefaultSchedulerContext.Send(slf.GetPID(), &netListenEvent{})

	return nil
}

//Shutdown doc
//@Summary Termination of service
//@Method Shutdown
func (slf *NetListener) Shutdown() {
	slf._opts.Sock.Close()
	slf.Service.Shutdown()
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

//Stopping Turn off network monitoring service
func (slf *NetListener) Stopping(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	slf.LogDebug("Listen Service Stopping")
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
	slf.LogDebug("Listen Service Stoped")
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

	if slf._opts.AsyncClose != nil {
		if err := slf._opts.AsyncClose(sockHandle); err != nil {
			slf.LogError("closed client Notification %+v", err)
		}
	}
	slf.LogDebug("closed client: %+v", sockHandle)
}

/*
//INetListenerDeleate Network listening commission
type INetListenerDeleate interface {
	Handshake(c INetClient) error
	Decode(context actor.Context, nets *NetListener, c INetClient) error
	UnOnlineNotification(h uint64) error
}

// NetListenService Network monitoring service
type NetListener struct {
	handler.Service

	NetListen  net.INetListen
	NetClients INetClientManager
	NetDeleate INetListenerDeleate
	NetMethod  NetMethodDispatch

	Addr       string //listening address
	CCMax      int    //Connector pipe buffer to small
	MaxClient  int
	ClientKeep uint64
}

//Initial Initialize the network listening service
func (slf *NetListener) Initial() {
Service.Initial()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&network.NetAccept{}, slf.OnAccept)
	slf.RegisterMethod(&network.NetChunk{}, slf.OnRecv)
	slf.RegisterMethod(&network.NetClose{}, slf.OnClose)
}

func (slf *NetListener) getDesc() string {
	return fmt.Sprintf("Network Listen [%s] ", slf.NetListen.Name())
}

//Started Turn on network monitoring service
func (slf *NetListener) Started(context actor.Context, sender *actor.PID, message interface{}) {
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
func (slf *NetListener) Stopping(context actor.Context,
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
func (slf *NetListener) OnAccept(context actor.Context,
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
func (slf *NetListener) OnRecv(context actor.Context,
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
			err = slf.NetDeleate.Decode(context, slf, c)
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
func (slf *NetListener) OnClose(context actor.Context,
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
func (slf *NetListener) Shutdown() {
	if slf.NetListen != nil {
		slf.NetListen.Close()
	}

	slf.Service.Shutdown()
}

//LogInfo Log information
func (slf *NetListener) LogInfo(frmt string, args ...interface{}) {
	slf.Service.LogInfo(slf.getDesc()+frmt, args...)
}

//LogError Record error log information
func (slf *NetListener) LogError(frmt string, args ...interface{}) {
	slf.Service.LogError(slf.getDesc()+frmt, args...)
}

//LogDebug Record debug log information
func (slf *NetListener) LogDebug(frmt string, args ...interface{}) {
	slf.Service.LogDebug(slf.getDesc()+frmt, args...)
}

//LogTrace Record trace log information
func (slf *NetListener) LogTrace(frmt string, args ...interface{}) {
	slf.Service.LogTrace(slf.getDesc()+frmt, args...)
}

//LogWarning Record warning log information
func (slf *NetListener) LogWarning(frmt string, args ...interface{}) {
	slf.Service.LogWarning(slf.getDesc()+frmt, args...)
}
*/
