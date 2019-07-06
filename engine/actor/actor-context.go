package actor

import (
  "time"
  "magicNet/engine/logger"
  "github.com/emirpasic/gods/stacks/linkedliststack"
)

type contextState int32

const (
  stateNone           contextState = iota
  stateAlive
  stateRestarting
  stateStopping
  stateStopped
)

func newActorContext(agnet *Agnets) *actorContext {
  this := &actorContext {
    agnet: agnet,
  }
  this.initActor()
  return this
}

type actorContext struct {
  actor               Actor
  agnet              *Agnets
  watchers            PIDSet
  self                *PID
  currentMessage      interface{}
  stash               *linkedliststack.Stack
  state               contextState
}

func (ctx *actorContext) initActor() {
  ctx.state = stateAlive
  ctx.actor = ctx.agnet.producer()
}

func (ctx *actorContext) Self() *PID {
	return ctx.self
}

func (ctx *actorContext) Sender() *PID {
	return UnWrapPackSender(ctx.currentMessage)
}

func (ctx *actorContext) Actor() Actor {
	return ctx.actor
}

func (ctx *actorContext) Message() interface{} {
	return UnWrapPackMessage(ctx.currentMessage)
}

func (ctx *actorContext) MessageHeader() ReadOnlyMessageHeader {
	return UnWrapPackHeader(ctx.currentMessage)
}

func (ctx *actorContext) Send(pid *PID, message interface{}) {
	ctx.sendUsrMessage(pid, message)
}

func (ctx *actorContext) sendUsrMessage(pid *PID, message interface{}) {
  pid.sendUsrMessage(message)
}

func (ctx *actorContext) Request(pid *PID, message interface{}) {
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  ctx.Self(),
	}

	ctx.sendUsrMessage(pid, e)
}

func (ctx *actorContext) RequestWithCustomSender(pid *PID, message interface{}, sender *PID) {
	env := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	ctx.sendUsrMessage(pid, env)
}

func (ctx *actorContext) RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	env := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  future.PID(),
	}
	ctx.sendUsrMessage(pid, env)

	return future
}

func (ctx *actorContext) Respond(response interface{}) {
	// If the message is addressed to nil forward it to the dead letter channel
	if ctx.Sender() == nil {
		deathLetter.SendUsrMessage(nil, response)
		return
	}

	ctx.Send(ctx.Sender(), response)
}


func (ctx *actorContext) Stash() {
	if ctx.stash == nil {
		ctx.stash = linkedliststack.New()
	}
	ctx.stash.Push(ctx.Message())
}

func (ctx *actorContext) Watch(who *PID) {
  who.sendSysMessage(&Watch{
    Watcher: ctx.self,
  })
}

func (ctx *actorContext) Unwatch(who *PID) {
  who.sendSysMessage(&Unwatch{
    Watcher: ctx.self,
  })
}

func (ctx *actorContext) Forward(pid *PID) {
	if msg, ok := ctx.currentMessage.(SystemMessage); ok {
    logger.Error(ctx.self.Id, "system message cannot be forwarded %v", msg)
		return
	}
	ctx.sendUsrMessage(pid, ctx.currentMessage)
}

func (ctx *actorContext) AwaitFuture(f *Future, cont func(res interface{}, err error)) {
	wrapper := func() {
		cont(f.result, f.err)
	}

	message := ctx.currentMessage
	// invoke the callback when the future completes
	f.continueWith(func(res interface{}, err error) {
		// send the wrapped callback as a continuation message to self
		ctx.self.sendSysMessage(&continuation{
			f:       wrapper,
			message: message,
		})
	})
}


func (ctx *actorContext) watch(watcher *PID) {
  ctx.watchers.Add(watcher)
}

func (ctx *actorContext) unwatch(watcher *PID) {
  ctx.watchers.Remove(watcher)
}


func (ctx *actorContext) InvokeUsrMessage(message interface{}) {
  if ctx.state == stateStopped {
    return
  }

  ctx.processMessage(message)
}

func (ctx *actorContext) processMessage(m interface{}) {
  ctx.currentMessage = m
  ctx.defaultReceive()
  ctx.currentMessage = nil
}

func (ctx *actorContext) Receive(pack *MessagePack) {
  ctx.currentMessage = pack
  ctx.defaultReceive()
  ctx.currentMessage = nil
}

func (ctx *actorContext) defaultReceive() {
  if _, ok := ctx.Message().(*Kill); ok {
    ctx.Stop(ctx.self)
    return
  }
  ctx.actor.Receive(Context(ctx))
}

func (ctx *actorContext) InvokeSysMessage(message interface{}) {
    switch msg := message.(type) {
    case *continuation:
      ctx.currentMessage = msg.message
      msg.f()
      ctx.currentMessage = nil
    case *Started:
      ctx.InvokeUsrMessage(msg)
    case *Watch:
      ctx.handleWatch(msg)
    case *Unwatch:
      ctx.handleUnWatch(msg)
    case *Stop:
      ctx.handleStop(msg)
    case *Terminated:
      ctx.handleTerminated(msg)
    default:
      logger.Error(ctx.self.Id, "unknown system message %v", msg)
    }
}

func (ctx *actorContext) handleWatch(msg *Watch) {
  if ctx.state >= stateStopping {
    msg.Watcher.sendSysMessage(&Terminated{
      Who: ctx.self,
    })
  } else {
    ctx.watch(msg.Watcher)
  }
}

func (ctx *actorContext) handleUnWatch(msg *Unwatch) {
    ctx.unwatch(msg.Watcher)
}

func (ctx *actorContext) handleStop(msg *Stop) {
  if ctx.state >= stateStopping {
    return
  }

  ctx.state = stateStopping
  ctx.InvokeUsrMessage(stoppingMessage)
  ctx.tryTerminate()
}

func (ctx *actorContext) handleTerminated(msg *Terminated) {
	ctx.InvokeUsrMessage(msg)
	ctx.tryTerminate()
}


func (ctx *actorContext) EscalateFailure(reason interface{}, message interface{}) {
  //TODO 考虑出错的问题
}

func (ctx *actorContext) tryTerminate() {
  if ctx.state == stateStopped {
    ctx.finalizeStop()
  }
}

func (ctx *actorContext) finalizeStop() {
	GlobalRegistry.UnRegister(ctx.self)
	ctx.InvokeUsrMessage(stoppedMessage)
	otherStopped := &Terminated{Who: ctx.self}
	// Notify watchers
  ctx.watchers.ForEach(func(i int, pid PID) {
      pid.sendSysMessage(otherStopped)
  })

	ctx.state = stateStopped
}

func (ctx *actorContext) Stop(pid *PID) {
  pid.ref().Stop(pid)
}

func (ctx *actorContext) StopFuture(pid *PID) *Future {
  future := NewFuture(10 * time.Second)

  pid.sendSysMessage(&Watch{Watcher: future.pid})
  ctx.Stop(pid)

  return future
}

func (ctx *actorContext) Kill(pid *PID) {
  pid.sendUsrMessage(&Kill{})
}

func (ctx *actorContext) KillFuture(pid *PID) *Future {
  future := NewFuture(10 * time.Second)

  pid.sendSysMessage(&Watch{Watcher: future.pid})
  ctx.Kill(pid)
  return future
}


func (ctx *actorContext) GoString() string {
	return ctx.self.String()
}

func (ctx *actorContext) String() string {
	return ctx.self.String()
}