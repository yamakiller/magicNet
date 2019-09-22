package actor

import (
	"time"
)

// SchedulerContext : 调度器
type SchedulerContext struct {
	headers messageHeader
}

// DefaultSchedulerContext :  默认推送Context
var DefaultSchedulerContext = &SchedulerContext{
	headers: DefaultMessageHeader,
}

// Sender : 推送者[无效]
func (sc *SchedulerContext) Sender() *PID {
	return nil
}

// Self : 自己[无效]
func (sc *SchedulerContext) Self() *PID {
	return nil
}

// Actor :  自己的Actor对象[无效]
func (sc *SchedulerContext) Actor() Actor {
	return nil
}

// Message : Message block [invalid]
func (sc *SchedulerContext) Message() interface{} {
	return nil
}

// MessageHeader : Header
func (sc *SchedulerContext) MessageHeader() ReadOnlyMessageHeader {
	return sc.headers
}

//SetHeaders : Set the basic message header
func (sc *SchedulerContext) SetHeaders(headers map[string]string) *SchedulerContext {
	sc.headers = headers
	return sc
}

// Send : Sending message to actor
func (sc *SchedulerContext) Send(pid *PID, message interface{}) {
	sc.sendUsrMessage(pid, message)
}

// Request : Request message
func (sc *SchedulerContext) Request(pid *PID, message interface{}) {
	sc.sendUsrMessage(pid, message)
}

// RequestWithCustomSender : Request a custom sender
func (sc *SchedulerContext) RequestWithCustomSender(pid *PID, message interface{}, sender *PID) {
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	sc.sendUsrMessage(pid, e)
}

// RequestFuture : Request pending reply
func (sc *SchedulerContext) RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  future.PID(),
	}
	sc.sendUsrMessage(pid, e)
	return future
}

// sendUsrMessage : Send a message
func (sc *SchedulerContext) sendUsrMessage(pid *PID, message interface{}) {
	pid.sendUsrMessage(message)
}

// Make : Maker
func (sc *SchedulerContext) Make(agnet *Agnets) *PID {
	pid, err := sc.MakeNamed(agnet, "")
	if err != nil {
		panic(err)
	}
	return pid
}

// MakeNamed : Maker with name
func (sc *SchedulerContext) MakeNamed(agnet *Agnets, name string) (*PID, error) {
	return agnet.make()
}

// Stop : Send stop message
func (sc *SchedulerContext) Stop(pid *PID) {
	pid.ref().Stop(pid)
}

// StopFuture : Send a stop message and wait for a reply
func (sc *SchedulerContext) StopFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	sc.Stop(pid)

	return future
}

// Kill : Kill Actor
func (sc *SchedulerContext) Kill(pid *PID) {
	pid.sendUsrMessage(&Kill{})
}

// KillFuture : Kill the Actor and wait for a reply
func (sc *SchedulerContext) KillFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	sc.Kill(pid)

	return future
}
