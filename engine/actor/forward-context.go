package actor

import (
  "time"
)


type ForwardContext struct {
  senderMiddleware SenderFunc
  headers          messageHeader
}

var DefaultForwardContext = &ForwardContext{
        nil,
        DefaultMessageHeader,
}

func (rf *ForwardContext) Sender() *PID {
	return nil
}

func (rf *ForwardContext) Self() *PID {
	return nil
}

func (rf *ForwardContext) Actor() Actor {
	return nil
}

func (fc *ForwardContext) Message()interface{} {
  return nil
}

func (fc *ForwardContext) MessageHeader() ReadOnlyMessageHeader {
  return fc.headers
}

func (fc *ForwardContext) Send(pid *PID, message interface{}) {
  fc.sendUsrMessage(pid, message)
}

func (fc *ForwardContext) Request(pid *PID, message interface{}) {
  fc.sendUsrMessage(pid, message)
}

func (fc *ForwardContext) RequestWithCustomSender(pid *PID, message interface{}, sender *PID) {
  e := &MessagePack {
    Header  : nil,
    Message : message,
    Sender  : sender,
  }
  fc.sendUsrMessage(pid, e)
}

func (fc *ForwardContext) RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
  future := NewFuture(timeout)
  e := &MessagePack {
    Header:  nil,
    Message: message,
    Sender:  future.PID(),
  }
  fc.sendUsrMessage(pid, e)
  return future
}

func (fc *ForwardContext) sendUsrMessage(pid *PID, message interface{}) {
  if fc.senderMiddleware != nil {
    fc.senderMiddleware(fc, pid, WrapPack(message))
  } else {
    pid.sendUsrMessage(message)
  }
}

func (fc *ForwardContext) Stop(pid *PID) {
  pid.ref().Stop(pid)
}

func (fc *ForwardContext) StopFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	fc.Stop(pid)

	return future
}

func (fc *ForwardContext) Kill(pid *PID ){
  pid.sendUsrMessage(&Kill{})
}

func (fc *ForwardContext) KillFuture(pid *PID) *Future {
  future := NewFuture(10 * time.Second)

  pid.sendSysMessage(&Watch{Watcher: future.pid})
  fc.Kill(pid)

  return future 
}
