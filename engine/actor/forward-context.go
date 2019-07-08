package actor

import (
	"time"
)

// ForwardContext : 基础推送Context
type ForwardContext struct {
	senderMiddleware SenderFunc
	headers          messageHeader
}

// DefaultForwardContext :  默认推送Context
var DefaultForwardContext = &ForwardContext{
	nil,
	DefaultMessageHeader,
}

// Sender : 推送者[无效]
func (fc *ForwardContext) Sender() *PID {
	return nil
}

// Self : 自己[无效]
func (fc *ForwardContext) Self() *PID {
	return nil
}

// Actor :  自己的Actor对象[无效]
func (fc *ForwardContext) Actor() Actor {
	return nil
}

// Message : 消息块[无效]
func (fc *ForwardContext) Message() interface{} {
	return nil
}

// MessageHeader : 消息头
func (fc *ForwardContext) MessageHeader() ReadOnlyMessageHeader {
	return fc.headers
}

// Send : 发送消息
func (fc *ForwardContext) Send(pid *PID, message interface{}) {
	fc.sendUsrMessage(pid, message)
}

// Request : 请求消息
func (fc *ForwardContext) Request(pid *PID, message interface{}) {
	fc.sendUsrMessage(pid, message)
}

// RequestWithCustomSender : 请求自定义发件人
func (fc *ForwardContext) RequestWithCustomSender(pid *PID, message interface{}, sender *PID) {
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  sender,
	}
	fc.sendUsrMessage(pid, e)
}

// RequestFuture : 请求等待回复
func (fc *ForwardContext) RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future {
	future := NewFuture(timeout)
	e := &MessagePack{
		Header:  nil,
		Message: message,
		Sender:  future.PID(),
	}
	fc.sendUsrMessage(pid, e)
	return future
}

// sendUsrMessage : 发送消息
func (fc *ForwardContext) sendUsrMessage(pid *PID, message interface{}) {
	if fc.senderMiddleware != nil {
		fc.senderMiddleware(fc, pid, WrapPack(message))
	} else {
		pid.sendUsrMessage(message)
	}
}

// Stop : 发送停止消息
func (fc *ForwardContext) Stop(pid *PID) {
	pid.ref().Stop(pid)
}

// StopFuture : 发送停止消息，并等待回复
func (fc *ForwardContext) StopFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	fc.Stop(pid)

	return future
}

// Kill : 杀死 Actor
func (fc *ForwardContext) Kill(pid *PID) {
	pid.sendUsrMessage(&Kill{})
}

// KillFuture : 杀死 Actor 并等待回复
func (fc *ForwardContext) KillFuture(pid *PID) *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	fc.Kill(pid)

	return future
}
