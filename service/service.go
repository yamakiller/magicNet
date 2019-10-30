package service

import (
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/util"
)

// MethodFunc : Service method function
type MethodFunc func(self actor.Context, sender *actor.PID, message interface{})

// IService  Service base class interface
type IService interface {
	actor.Actor

	GetPID() *actor.PID
	Name() string
	Key() string
	ID() uint32

	Init()
	Assignment(context actor.Context)
	Shutdown()

	RegisterMethod(key interface{}, method MethodFunc)

	withName(n string)

	withWait(wait *sync.WaitGroup)

	LogInfo(frmt string, args ...interface{})
	LogError(frmt string, args ...interface{})
	LogDebug(frmt string, args ...interface{})
	LogTrace(frmt string, args ...interface{})
	LogWarning(frmt string, args ...interface{})
}

// Service server base class
type Service struct {
	pid    *actor.PID
	name   string
	wait   *sync.WaitGroup
	method map[interface{}]MethodFunc
}

// Init : Initialization service
func (slf *Service) Init() {
	slf.method = make(map[interface{}]MethodFunc, 16)
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
	slf.RegisterMethod(&actor.Terminated{}, slf.Terminated)
}

// Receive : Receive message
func (slf *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.MessagePack:
		f, ok := slf.method[reflect.TypeOf(msg.Message)]
		if ok {
			f(context, msg.Sender, msg.Message)
		}
		logger.Error(context.Self().ID, "service unknown message:%+v,sender:%+v", msg, msg.Sender)
	default:
		f, ok := slf.method[reflect.TypeOf(msg)]
		if ok {
			f(context, nil, msg)
			break
		}
		logger.Error(context.Self().ID, "service unknown message:%+v", msg)
	}
}

//Assignment Service initial value association
func (slf *Service) Assignment(context actor.Context) {
	slf.pid = context.Self()
	slf.name = slf.name + "$" + strconv.Itoa(int(slf.pid.ID))
}

// Started : Service start function
func (slf *Service) Started(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	if slf.pid == nil {
		slf.Assignment(context)
	}
	if slf.wait != nil {
		slf.wait.Done()
	}
}

//Stopping : Service stopped
func (slf *Service) Stopping(context actor.Context, sender *actor.PID, message interface{}) {
}

// Stoped : Service stop closing handler
func (slf *Service) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	for k := range slf.method {
		delete(slf.method, k)
	}

	if slf.wait != nil {
		slf.wait.Done()
	}
}

// Terminated : The service is terminated and can be destroyed
func (slf *Service) Terminated(context actor.Context, sender *actor.PID, message interface{}) {
	slf.Shutdown()
}

// Shutdown : Proactively shut down the service
func (slf *Service) Shutdown() {
	if slf.pid == nil {
		return
	}
	slf.pid.Stop()
	slf.wait.Wait()
}

// Name : Get the name of the service
func (slf *Service) Name() string {
	return slf.name
}

// Key : Returns the Key name of the service
func (slf *Service) Key() string {
	ix := strings.IndexByte(slf.name, '$')
	if ix <= 0 {
		return slf.name
	}

	return util.SubStr2(slf.name, 0, ix)
}

//GetPID Return the pid object
func (slf *Service) GetPID() *actor.PID {
	return slf.pid
}

// ID Returns the unique number of the service
func (slf *Service) ID() uint32 {
	return slf.pid.ID
}

// RegisterMethod : Registration (convention/agreement) method
func (slf *Service) RegisterMethod(key interface{}, method MethodFunc) {
	slf.method[reflect.TypeOf(key)] = method
}

func (slf *Service) withName(n string) {
	slf.name = n
}

func (slf *Service) withWait(wait *sync.WaitGroup) {
	slf.wait = wait
}

//LogInfo Log information
func (slf *Service) LogInfo(frmt string, args ...interface{}) {
	logger.Info(slf.ID(), frmt, args...)
}

//LogError Record error log information
func (slf *Service) LogError(frmt string, args ...interface{}) {
	logger.Error(slf.ID(), frmt, args...)
}

//LogDebug Record debug log information
func (slf *Service) LogDebug(frmt string, args ...interface{}) {
	logger.Debug(slf.ID(), frmt, args...)
}

//LogTrace Record trace log information
func (slf *Service) LogTrace(frmt string, args ...interface{}) {
	logger.Trace(slf.ID(), frmt, args...)
}

//LogWarning Record warning log information
func (slf *Service) LogWarning(frmt string, args ...interface{}) {
	logger.Warning(slf.ID(), frmt, args...)
}

// Make : Service creator
func Make(name string, f func() IService) IService {
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
