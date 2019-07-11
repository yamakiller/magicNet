package service

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"strings"
	"sync"
)

// MethodFunc : 服务方法函数
type MethodFunc func(self actor.Context, message interface{})

// IService 服务基础类接口
type IService interface {
	actor.Actor

	Name() string
	Key() string
	ID() uint32

	Started(context actor.Context)
	Stoped(context actor.Context)
	Terminated(context actor.Context)
	Shutdown()

	RegisterMethod(key interface{}, method MethodFunc)

	setName(n string)
}

// Service 服务基类
type Service struct {
	pid    *actor.PID
	name   string
	wait   sync.WaitGroup
	method map[interface{}]MethodFunc
}

// Receive : 接收消息
func (srv *Service) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		srv.Started(context)
		break
	case *actor.Stopped:
		srv.Stoped(context)
		break
	case *actor.Terminated:
		srv.Terminated(context)
	default:
		f, ok := srv.method[msg]
		if ok {
			f(context, msg)
			break
		}
		logger.Error(context.Self().ID, "service unknown message:%v", msg)
	}
}

// Started : 服务的启动函数
func (srv *Service) Started(context actor.Context) {

}

// Stoped : 服务停止收尾处理函数
func (srv *Service) Stoped(context actor.Context) {
	for k := range srv.method {
		delete(srv.method, k)
	}
}

// Terminated : 服务被终止可以被销毁
func (srv *Service) Terminated(context actor.Context) {
	srv.wait.Done()
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

// ID : 获取服务的唯一编号
func (srv *Service) ID() uint32 {
	return srv.pid.ID
}

// RegisterMethod : 注册(约定/协议)方法
func (srv *Service) RegisterMethod(key interface{}, method MethodFunc) {
	srv.method[key] = method
}

func (srv *Service) setName(n string) {
	srv.name = n
}

var serviceAgnet = Agnets{}

// Make : 服务制作器
func Make(name string, f func() IService) IService {
	// 需要调整

	/*schedulerContext := actor.DefaultSchedulerContext
	srv := f()
	srv.setName(name)

	serivceAgnet := actor.Agnets{maker: func(agnet *Agnets) (*PID, error) {
		ctx := newActorContext(agnet)
		mb := agnet.produceMailbox()
		dp := agnet.getDispatcher()
		proc := NewActorProcess(mb)
		pid := &PID{}
		globalRegistry.Register(pid, proc)
		ctx.self = pid
		mb.Start()
		mb.RegisterHandlers(ctx, dp)
		mb.PostSysMessage(startedMessage)
		return pid, nil
	}}

	serivceAgnet.SetMakeActor(func() actor.Actor {
		return srv
	})

	schedulerContext.Make(&serivceAgnet)
	return srv*/

	/*ref := &process{}
	proxy, absent := actor.ProcessRegistry.Add(ref, id)
	if !absent {
		return proxy, actor.ErrNameExists
	}

	var pc = *props
	pc.WithSpawnFunc(nil)
	ref.state = config.CreateRouterState()

	if config.RouterType() == GroupRouterType {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ref.router, _ = actor.DefaultSpawner(id+"/router", actor.PropsFromProducer(func() actor.Actor {
			return &groupRouterActor{
				props:  &pc,
				config: config,
				state:  ref.state,
				wg:     wg,
			}
		}), parentContext)
		wg.Wait() // wait for routerActor to start
	} else {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ref.router, _ = actor.DefaultSpawner(id+"/router", actor.PropsFromProducer(func() actor.Actor {
			return &poolRouterActor{
				props:  &pc,
				config: config,
				state:  ref.state,
				wg:     wg,
			}
		}), parentContext)
		wg.Wait() // wait for routerActor to start
	}

	ref.parent = parentContext.Self()
	return proxy, nil*/
	return nil
}
