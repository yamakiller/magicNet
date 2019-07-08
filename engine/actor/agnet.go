package actor

import (
	"magicNet/engine/mailbox"
	"unsafe"
)

// SpawnFunc : 创建ActorContext函数
type SpawnFunc func(agnet *Agnets) (*PID, error)

var (
	defaultDispatcher      = mailbox.NewGoroutineDispatcher(300)
	defaultMailboxProducer = mailbox.Unbounded()
	defaultSpawner         = func(agnet *Agnets) (*PID, error) {
		ctx := newActorContext(agnet)
		mb := agnet.produceMailbox()
		dp := agnet.getDispatcher()
		proc := NewActorProcess(mb)
		pid := &PID{p: (*Process)(unsafe.Pointer(proc))}
		GlobalRegistry.Register(pid)
		ctx.self = pid
		mb.Start()
		mb.RegisterHandlers(ctx, dp)
		mb.PostSysMessage(startedMessage)
		return pid, nil
	}
)

// DefaultSpawner : 默认创建函数代理
var DefaultSpawner SpawnFunc = defaultSpawner

// Agnets : ActorContext 代理对象
type Agnets struct {
	spawner         SpawnFunc
	newactor        NewActor
	mailboxProducer mailbox.Producer
	dispatcher      mailbox.Dispatcher
}

func (agnet *Agnets) getSpawner() SpawnFunc {
	if agnet.spawner == nil {
		return defaultSpawner
	}
	return agnet.spawner
}

func (agnet *Agnets) getDispatcher() mailbox.Dispatcher {
	if agnet.dispatcher == nil {
		return defaultDispatcher
	}
	return agnet.dispatcher
}

func (agnet *Agnets) produceMailbox() mailbox.Mailbox {
	if agnet.mailboxProducer == nil {
		return defaultMailboxProducer()
	}
	return agnet.mailboxProducer()
}
