package actor

import (
  "unsafe"
  "magicNet/engine/mailbox"
)


// 默认值
type SpawnFunc func(agnet *Agnets) (*PID, error)

var (
  defaultDispatcher      = mailbox.NewGoroutineDispatcher(300)
  defaultMailboxProducer = mailbox.Unbounded()
  defaultSpawner         = func(agnet *Agnets) (*PID, error) {
      ctx := newActorContext(agnet)
      mb  := agnet.produceMailbox()
      dp  := agnet.getDispatcher()
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

var DefaultSpawner SpawnFunc = defaultSpawner

type Agnets struct {
  spawner           SpawnFunc
  producer          Producer
  mailboxProducer   mailbox.Producer
  dispatcher        mailbox.Dispatcher
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
