package actor

import (
	"magicNet/engine/eventchannel"
	"magicNet/engine/logger"
	"magicNet/engine/util"
)

type deathLetterProcess struct{}

var (
	deathLetter           Process = &deathLetterProcess{}
	deathLetterSubscriber *eventchannel.Subscription
)

func init() {
	deathLetterSubscriber = eventchannel.Subscribe(func(evt interface{}) {
		if deathLetter, ok := evt.(*DeadLetterEvent); ok {
			util.Assert(deathLetter.Sender != nil && deathLetter.PID != nil, "deathLetter sender or pid is nil")
			logger.Debug(deathLetter.Sender.GetId(), "DeathLetter Dest PID :%s", deathLetter.PID.String())
		}
	})

	eventchannel.Subscribe(func(evt interface{}) {
		if deathLetter, ok := evt.(*DeadLetterEvent); ok {
			if m, ok := deathLetter.Message.(*Watch); ok {
				m.Watcher.sendSysMessage(&Terminated{AddressTerminated: false, Who: deathLetter.PID})
			}
		}
	})
}

// DeadLetterEvent : 死亡消息
type DeadLetterEvent struct {
	PID     *PID        // The invalid process, to which the message was sent
	Message interface{} // The message that could not be delivered
	Sender  *PID        // the process that sent the Message
}

// SendUsrMessage ： 发送死亡消息
func (*deathLetterProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	eventchannel.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	})
}

// SendSysMessage : 发送死亡消息
func (*deathLetterProcess) SendSysMessage(pid *PID, message interface{}) {
	eventchannel.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: message,
	})
}

// Stop: 发送停止消息
func (ref *deathLetterProcess) Stop(pid *PID) {
	ref.SendSysMessage(pid, stopMessage)
}
