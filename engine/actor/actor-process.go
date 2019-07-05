package actor

import (
  "sync/atomic"
  "magicNet/engine/mailbox"
)

type ActorProcess struct {
  mailbox mailbox.Mailbox
  death   int32
}

func NewActorProcess(mailbox mailbox.Mailbox) *ActorProcess {
  return &ActorProcess{mailbox: mailbox}
}

func (a *ActorProcess) SendUsrMessage(pid *PID, message interface{}) {
  a.mailbox.PostUsrMessage(message)
}

func (a *ActorProcess) SendSysMessage(pid *PID, message interface{}) {
  a.mailbox.PostSysMessage(message)
}

func (a *ActorProcess) Stop(pid *PID) {
  atomic.StoreInt32(&a.death, 1)
  a.SendSysMessage(pid, stopMessage)
}
