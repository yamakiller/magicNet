package actor

import (
	"magicNet/engine/mailbox"
	"sync/atomic"
)

// AtrProcess : Actor 处理模块
type AtrProcess struct {
	mailbox mailbox.Mailbox
	death   int32
}

// NewActorProcess : 创建一个 ActorProcess
func NewActorProcess(mailbox mailbox.Mailbox) *AtrProcess {
	return &AtrProcess{mailbox: mailbox}
}

// SendUsrMessage : 发送用户级消息
func (a *AtrProcess) SendUsrMessage(pid *PID, message interface{}) {
	a.mailbox.PostUsrMessage(message)
}

// SendSysMessage : 发送系统级消息
func (a *AtrProcess) SendSysMessage(pid *PID, message interface{}) {
	a.mailbox.PostSysMessage(message)
}

// Stop : 发送停止Actor消息
func (a *AtrProcess) Stop(pid *PID) {
	atomic.StoreInt32(&a.death, 1)
	a.SendSysMessage(pid, stopMessage)
}
