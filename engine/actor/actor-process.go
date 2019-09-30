package actor

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月05日 18:00:31
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月21日 10:19:54
 * @Description:Actor 处理消息处理模块
 */

import (
	"sync/atomic"

	"github.com/yamakiller/magicNet/engine/mailbox"
)

// AtrProcess Actor Processing module
type AtrProcess struct {
	mailbox mailbox.Mailbox
	death   int32
}

// NewActorProcess Create an Actor processor
func NewActorProcess(mailbox mailbox.Mailbox) *AtrProcess {
	return &AtrProcess{mailbox: mailbox}
}

// SendUsrMessage Send user level messages
func (a *AtrProcess) SendUsrMessage(pid *PID, message interface{}) {
	a.mailbox.PostUsrMessage(message)
}

// SendSysMessage Send system level messages
func (a *AtrProcess) SendSysMessage(pid *PID, message interface{}) {
	a.mailbox.PostSysMessage(message)
}

// OverloadUsrMessage A warning to detect if a user's mailbox has been expanded
func (a *AtrProcess) OverloadUsrMessage() int {
	return a.mailbox.OverloadUsrMessage()
}

// Stop Send stop Actor message
func (a *AtrProcess) Stop(pid *PID) {
	atomic.StoreInt32(&a.death, 1)
	a.SendSysMessage(pid, stopMessage)
}
