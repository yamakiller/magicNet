package actor

import (
	"time"
)

// SchedulerContext : 调度器
type SchedulerContext struct {
	senderMiddleware SenderFunc
	headers          messageHeader
}

// DefaultSchedulerContext :  默认推送Context
var DefaultSchedulerContext = &SchedulerContext{
	nil,
	DefaultMessageHeader,
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

// Message : 消息块[无效]
func (sc *SchedulerContext) Message() interface{} {
	return nil
}

// MessageHeader : 消息头
func (sc *SchedulerContext) MessageHeader() ReadOnlyMessageHeader {
	return sc.headers
}

//SetHeaders : 设置基础消息头
func (sc *SchedulerContext) SetHeaders(headers map[string]string) *SchedulerContext {
	sc.headers = headers
	return sc
}

// Send : 发送消息
func (sc *SchedulerContext) Send(pid *PID, message interface{}) {
	sc.sendUsrMessage(pid, message)
}

// Request : 请求消息
func (sc *SchedulerContext) Request(pid *PID, message interface{}) {
	sc.sendUsrMessage(pid, message)
}

// RequestWithCustomSender : 请求自定义发件人
func (sc *SchedulerContext) RequestWithCustomSender(pid *PID, message interface{}, sender *PID) {
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	sc.sendUsrMessage(pid, e)
}

// RequestFuture : 请求等待回复
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

// sendUsrMessage : 发送消息
func (sc *SchedulerContext) sendUsrMessage(pid *PID, message interface{}) {
	if sc.senderMiddleware != nil {
		sc.senderMiddleware(sc, pid, WrapPack(message))
	} else {
		pid.sendUsrMessage(message)
	}
}

// Make : 制作器
func (sc *SchedulerContext) Make(agnet *Agnets) *PID {
	pid, err := sc.MakeNamed(agnet, "")
	if err != nil {
		panic(err)
	}
	return pid
}

// MakeNamed : 带名字的制作器
func (sc *SchedulerContext) MakeNamed(agnet *Agnets, name string) (*PID, error) {
	return agnet.make()
}

// Stop : 发送停止消息
func (sc *SchedulerContext) Stop(pid *PID) {
	pid.ref().Stop(pid)
}

// StopFuture : 发送停止消息，并等待回复
func (sc *SchedulerContext) StopFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	sc.Stop(pid)

	return future
}

// Kill : 杀死 Actor
func (sc *SchedulerContext) Kill(pid *PID) {
	pid.sendUsrMessage(&Kill{})
}

// KillFuture : 杀死 Actor 并等待回复
func (sc *SchedulerContext) KillFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	sc.Kill(pid)

	return future
}
