package actor

import (
	"magicNet/engine/mailbox"
	"sync/atomic"
)

// ActorProcess : Actor 处理模块
type ActorProcess struct {
	mailbox mailbox.Mailbox
	death   int32
}

// NewActorProcess : 创建一个 ActorProcess
func NewActorProcess(mailbox mailbox.Mailbox) *ActorProcess {
	return &ActorProcess{mailbox: mailbox}
}

// SendUsrMessage : 发送用户级消息
func (a *ActorProcess) SendUsrMessage(pid *PID, message interface{}) {
	a.mailbox.PostUsrMessage(message)
}

// SendSysMessage : 发送系统级消息
func (a *ActorProcess) SendSysMessage(pid *PID, message interface{}) {
	a.mailbox.PostSysMessage(message)
}

// Stop : 发送停止Actor消息
func (a *ActorProcess) Stop(pid *PID) {
	atomic.StoreInt32(&a.death, 1)
	a.SendSysMessage(pid, stopMessage)
}
