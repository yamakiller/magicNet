package actor

import (
	"github.com/yamakiller/magicNet/engine/mailbox"
)

// MakeFunc : 创建ActorContext函数
type MakeFunc func(agnet *Agnets) (*PID, error)

var (
	defaultDispatcher  = mailbox.NewGoroutineDispatcher(300)
	defaultMailboxMake = mailbox.Unbounded()
	defaultMaker       = func(agnet *Agnets) (*PID, error) {
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
	}
)

// DefaultMaker : 默认创建函数代理
var DefaultMaker MakeFunc = defaultMaker

// Agnets : ActorContext 代理对象
type Agnets struct {
	maker       MakeFunc
	actorMake   MakeActor
	mailboxMake mailbox.Make
	dispatcher  mailbox.Dispatcher
}

func (agnet *Agnets) getMaker() MakeFunc {
	if agnet.maker == nil {
		return defaultMaker
	}
	return agnet.maker
}

func (agnet *Agnets) getDispatcher() mailbox.Dispatcher {
	if agnet.dispatcher == nil {
		return defaultDispatcher
	}
	return agnet.dispatcher
}

func (agnet *Agnets) produceMailbox() mailbox.Mailbox {
	if agnet.mailboxMake == nil {
		return defaultMailboxMake()
	}
	return agnet.mailboxMake()
}

// SetMakeFunc : 设置基础制造机器
func (agnet *Agnets) SetMakeFunc(maker MakeFunc) *Agnets {
	agnet.maker = maker
	return agnet
}

// SetMakeActor : 设置Actor创建器
func (agnet *Agnets) SetMakeActor(m MakeActor) *Agnets {
	agnet.actorMake = m
	return agnet
}

// SetDispatcher ：设置消息分发器
func (agnet *Agnets) SetDispatcher(dispatcher mailbox.Dispatcher) *Agnets {
	agnet.dispatcher = dispatcher
	return agnet
}

// SetMailboxMake : 设置邮箱创建器
func (agnet *Agnets) SetMailboxMake(m mailbox.Make) *Agnets {
	agnet.mailboxMake = m
	return agnet
}

func (agnet *Agnets) make() (*PID, error) {
	return agnet.getMaker()(agnet)
}

// AgnetFromActorMaker : 创建一个分配给Actor制造的代理
func AgnetFromActorMaker(maker MakeActor) *Agnets {
	return &Agnets{
		actorMake: maker,
	}
}

// AgnetFromFunc 使用指定为actor生成器给定receive 函数创建一个代理
func AgnetFromFunc(f AtrFunc) *Agnets {
	return AgnetFromActorMaker(func() Actor { return f })
}
