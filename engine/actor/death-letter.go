package actor

import (
	"magicNet/engine/evtchan"
	"magicNet/engine/logger"
	"magicNet/engine/util"
)

type deathLetterProcess struct{}

var (
	deathLetter           Process = &deathLetterProcess{}
	deathLetterSubscriber *evtchan.Subscription
)

func init() {
	deathLetterSubscriber = evtchan.Subscribe(func(evt interface{}) {
		if deathLetter, ok := evt.(*DeadLetterEvent); ok {
			util.Assert(deathLetter.Sender != nil && deathLetter.PID != nil, "deathLetter sender or pid is nil")
			logger.Debug(deathLetter.Sender.GetID(), "DeathLetter Dest PID :%s", deathLetter.PID.String())
		}
	})

	evtchan.Subscribe(func(evt interface{}) {
		if deathLetter, ok := evt.(*DeadLetterEvent); ok {
			if m, ok := deathLetter.Message.(*Watch); ok {
				m.Watcher.sendSysMessage(&Terminated{AddressTerminated: false, Who: deathLetter.PID})
			}
		}
	})
}

// DeadLetterEvent : 死亡消息
type DeadLetterEvent struct {
	PID     *PID
	Message interface{}
	Sender  *PID
}

// SendUsrMessage ： 发送死亡消息
func (*deathLetterProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	})
}

// SendSysMessage : 发送死亡消息
func (*deathLetterProcess) SendSysMessage(pid *PID, message interface{}) {
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: message,
	})
}

// Stop: 发送停止消息
func (ref *deathLetterProcess) Stop(pid *PID) {
	ref.SendSysMessage(pid, stopMessage)
}
