package service

import (
	"reflect"
	"strings"
	"sync"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/util"
)

// MethodFunc : Service method function
type MethodFunc func(self actor.Context, message interface{})

// IService  Service base class interface
type IService interface {
	actor.Actor

	Name() string
	Key() string
	ID() uint32

	Init()
	Shutdown()

	RegisterMethod(key interface{}, method MethodFunc)

	GetMethod(key interface{}) MethodFunc

	withName(n string)

	withWait(wait *sync.WaitGroup)

	LogInfo(frmt string, args ...interface{})
	LogError(frmt string, args ...interface{})
	LogDebug(frmt string, args ...interface{})
	LogTrace(frmt string, args ...interface{})
	LogWarning(frmt string, args ...interface{})
}

// Service 服务基类
type Service struct {
	pid    *actor.PID
	name   string
	wait   *sync.WaitGroup
	method map[interface{}]MethodFunc
}

// Init : 初始化服务
func (srv *Service) Init() {
	srv.method = make(map[interface{}]MethodFunc, 16)
	srv.RegisterMethod(&actor.Terminated{}, srv.Terminated)
}

// Receive : 接收消息
func (srv *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	default:
		f, ok := srv.method[reflect.TypeOf(msg)]
		if ok {
			f(context, msg)
			break
		}
		logger.Error(context.Self().ID, "service unknown message:%+v", msg)
	}
}

// Started : 服务的启动函数
func (srv *Service) Started(context actor.Context, message interface{}) {
	srv.pid = context.Self()
	srv.name = srv.name + "$" + srv.pid.String()
	if srv.wait != nil {
		srv.wait.Done()
	}
}

// Stoped : 服务停止收尾处理函数
func (srv *Service) Stoped(context actor.Context, message interface{}) {
	for k := range srv.method {
		delete(srv.method, k)
	}
}

// Terminated : 服务被终止可以被销毁
func (srv *Service) Terminated(context actor.Context, message interface{}) {
	if srv.wait != nil {
		srv.wait.Done()
	}
}

// Shutdown : 主动关闭服务
func (srv *Service) Shutdown() {
	if srv.pid == nil {
		return
	}
	srv.pid.Stop()
	srv.wait.Wait()
}

// Name : 获取服务的名字
func (srv *Service) Name() string {
	return srv.name
}

// Key : 获取服务的Key 名字
func (srv *Service) Key() string {
	ix := strings.IndexByte(srv.name, '$')
	if ix <= 0 {
		return srv.name
	}

	return util.SubStr2(srv.name, 0, ix)
}

// ID : returns the unique number of the service
func (srv *Service) ID() uint32 {
	return srv.pid.ID
}

// RegisterMethod : Registration (convention/agreement) method
func (srv *Service) RegisterMethod(key interface{}, method MethodFunc) {
	srv.method[reflect.TypeOf(key)] = method
}

//GetMethod Return the method corresponding to the protocol
func (srv *Service) GetMethod(key interface{}) MethodFunc {
	if r, ok := srv.method[key]; ok {
		return r
	}
	return nil
}

func (srv *Service) withName(n string) {
	srv.name = n
}

func (srv *Service) withWait(wait *sync.WaitGroup) {
	srv.wait = wait
}

func(srv *Service) LogInfo(frmt string, args ...interface{}) {
	logger.Info(srv.ID(), frmt, args)
}

func(srv *Service) LogError(frmt string, args ...interface{}) {
	logger.Error(srv.ID(), frmt, args)
}

func(srv *Service) LogDebug(frmt string, args ...interface{}) {
	logger.Debug(srv.ID(), frmt, args)
}


func(srv *Service) LogTrace(frmt string, args ...interface{}) {
	logger.Trace(srv.ID(), frmt, args)
}

func(srv *Service) LogWarning(frmt string, args ...interface{}) {
	logger.Warning(srv.ID(),frmt, args)
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
