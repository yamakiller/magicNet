package handler

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/engine/actor"
)

//MethodFunc desc
//@type MethodFunc desc: Service method function
//@param (actor.Context) a actor context
//@param (actor.PID) sender actor ID
//@param (interface{}) a message
type MethodFunc func(self actor.Context, sender *actor.PID, message interface{})

//IService desc
//@interface IService desc: Service base class interface
//@inherit (actor.Actor)
//@method (GetPID() *actor.PID ) return this id
//@method (Name() string) return this name
//@method (Key() string) return this pid is string
//@method (ID() uint32) return this pid=>id
//@method (Init()) initialization this
//@method (Shutdown) shutdown this service
//@method (RegisterMethod) register event call method
//@method (withPID(context actor.Context)) assignment this pid
//@method (withName)
//@method (widthWait)
//@method (LogInfo)
//@method (LogError)
//@method (LogDebug)
//@method (LogTrace)
//@method (LogWarning)
type IService interface {
	actor.Actor

	GetPID() *actor.PID
	Name() string
	Key() string
	ID() uint32

	Initial()
	Shutdown()
	RegisterMethod(key interface{}, method MethodFunc)

	//WithPID(context actor.Context)
	withName(n string)
	withWait(wait *sync.WaitGroup)

	LogInfo(frmt string, args ...interface{})
	LogError(frmt string, args ...interface{})
	LogDebug(frmt string, args ...interface{})
	LogTrace(frmt string, args ...interface{})
	LogWarning(frmt string, args ...interface{})
}

//Service desc
//@struct Service desc: server base class
//@member (*actor.PID) this id
//@member (string) this server name
//@member (*sync.WaitGroup)
//@member (map[interface{}]MethodFunc) event method map
type Service struct {
	_pid    *actor.PID
	_name   string
	_wait   *sync.WaitGroup
	_method map[interface{}]MethodFunc
}

//Initial desc
//@method Initial desc: Initialization service
func (slf *Service) Initial() {
	if slf._method == nil {
		slf._method = make(map[interface{}]MethodFunc)
	}

	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
	slf.RegisterMethod(&actor.Terminated{}, slf.Terminated)
}

//Receive desc
//@method Receive desc: Receive message and Scheduling
//@param (actor.Context) source actor context
func (slf *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.MessagePack:
		f, ok := slf._method[reflect.TypeOf(msg.Message)]
		if ok {
			slf.withPID(msg.Message, context)
			f(context, msg.Sender, msg.Message)
		}
		logger.Error(context.Self().ID, "service unknown message:%+v,sender:%+v", msg, msg.Sender)
	default:
		f, ok := slf._method[reflect.TypeOf(msg)]
		if ok {
			slf.withPID(msg, context)
			f(context, nil, msg)
			break
		}
		logger.Error(context.Self().ID, "service unknown message:%+v", msg)
	}
}

//withPID desc
//@method withPID desc: Service initial value association
//@param (actor.Context) this service context
func (slf *Service) withPID(msg interface{}, context actor.Context) {
	if slf._pid == nil {
		switch msg.(type) {
		case *actor.Started:
			slf._pid = context.Self()
			slf._name = slf._name + "$" + strconv.Itoa(int(slf._pid.ID))
		default:
		}
	}

}

//Started desc
//@method Started desc: Started Event Call Function
//@param (actor.Context) source actor context
//@param (*actor.PID) sender actor ID
//@param (interface{}) a message
func (slf *Service) Started(context actor.Context,
	sender *actor.PID,
	message interface{}) {
	if slf._wait != nil {
		slf._wait.Done()
	}
}

//Stopping desc
//@method Stopping desc: Stopping Event Call Function
//@param (actor.Context) source actor context
//@param (*actor.PID) sender actor ID
//@param (interface{}) a message
func (slf *Service) Stopping(context actor.Context, sender *actor.PID, message interface{}) {
}

//Stoped desc
//@method Stopping desc: Stopped Event Call Function
//@param (actor.Context) source actor context
//@param (*actor.PID) sender actor ID
//@param (interface{}) a message
func (slf *Service) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	for k := range slf._method {
		delete(slf._method, k)
	}

	slf._pid = nil
	slf._name = ""

	if slf._wait != nil {
		slf._wait.Done()
	}
}

//Terminated desc
//@method Terminated desc: Terminated Event Call Function
//@param (actor.Context) source actor context
//@param (*actor.PID) sender actor ID
//@param (interface{}) a message
func (slf *Service) Terminated(context actor.Context, sender *actor.PID, message interface{}) {
	//slf.Shutdown()
}

//Shutdown desc
//@method Shutdown desc: Shutdown service
func (slf *Service) Shutdown() {
	if slf._pid == nil {
		return
	}
	slf._pid.Stop()
	slf._wait.Wait()
}

//Name desc
//@method Name desc: Return the name of the service
//@return (string) name
func (slf *Service) Name() string {
	return slf._name
}

//Key desc
//@method Key desc: Returns the Key name of the service
//@return (string) pid=>key
func (slf *Service) Key() string {
	ix := strings.IndexByte(slf._name, '$')
	if ix <= 0 {
		return slf._name
	}

	return util.SubStr2(slf._name, 0, ix)
}

//GetPID desc
//@method GetPID desc: Return the pid object
//@return (*actor.PID) actor ID
func (slf *Service) GetPID() *actor.PID {
	return slf._pid
}

//ID desc
//@method ID desc: Returns this service pid=>id
//@return (uint32) ID
func (slf *Service) ID() uint32 {
	return slf._pid.ID
}

//RegisterMethod desc
//@method RegisterMethod desc: Registration (convention/agreement) method
//@param (interface{}) event map key
//@param (MethodFunc) Function object
func (slf *Service) RegisterMethod(key interface{}, method MethodFunc) {
	slf._method[reflect.TypeOf(key)] = method
}

func (slf *Service) withName(n string) {
	slf._name = n
}

func (slf *Service) withWait(wait *sync.WaitGroup) {
	slf._wait = wait
}

//LogInfo desc
//@method LogInfo desc: Log information
//@param  (string) format string
//@param  (...interface{}) format args
func (slf *Service) LogInfo(frmt string, args ...interface{}) {
	logger.Info(slf.ID(), frmt, args...)
}

//LogError desc
//@method LogError desc: Record error log information
//@param  (string) format string
//@param  (...interface{}) format args
func (slf *Service) LogError(frmt string, args ...interface{}) {
	logger.Error(slf.ID(), frmt, args...)
}

//LogDebug desc
//@method LogDebug desc: Record debug log information
//@param  (string) format string
//@param  (...interface{}) format args
func (slf *Service) LogDebug(frmt string, args ...interface{}) {
	logger.Debug(slf.ID(), frmt, args...)
}

//LogTrace desc
//@method LogTrace desc: Record trace log information
//@param  (string) format string
//@param  (...interface{}) format args
func (slf *Service) LogTrace(frmt string, args ...interface{}) {
	logger.Trace(slf.ID(), frmt, args...)
}

//LogWarning desc
//@method LogWarning desc: Record warning log information
//@param  (string) format string
//@param  (...interface{}) format args
func (slf *Service) LogWarning(frmt string, args ...interface{}) {
	logger.Warning(slf.ID(), frmt, args...)
}

//Spawn desc
//@method Spawn desc: Service creator function
//@param (string) service name
//@param (func() IService) service maker(function)
//@param (IService) service
func Spawn(name string, f func() IService) IService {
	wgn := &sync.WaitGroup{}
	srv := f()
	srv.withName(name)
	srv.withWait(wgn)
	wgn.Add(1)
	actor.DefaultMaker(actor.AgnetFromActorMaker(func() actor.Actor {
		return srv
	}))
	wgn.Wait()
	wgn.Add(1)

	return srv
}