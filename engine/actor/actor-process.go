package actor

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月05日 18:00:31
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 15:52:28
 * @Description:Actor 处理消息处理模块
 */

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

// OverloadUsrMessage : 检测用户邮箱是否有扩容的警告
func (a *AtrProcess) OverloadUsrMessage() int {
	return a.mailbox.OverloadUsrMessage()
}

// Stop : 发送停止Actor消息
func (a *AtrProcess) Stop(pid *PID) {
	atomic.StoreInt32(&a.death, 1)
	a.SendSysMessage(pid, stopMessage)
}
