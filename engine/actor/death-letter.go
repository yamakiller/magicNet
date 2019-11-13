package actor

import (
	"github.com/yamakiller/magicNet/engine/evtchan"
	"github.com/yamakiller/magicNet/engine/logger"
)

type deathLetterProcess struct{}

var (
	deathLetter           Process = &deathLetterProcess{}
	deathLetterSubscriber *evtchan.Subscription
)

func init() {
	deathLetterSubscriber = evtchan.Subscribe(func(evt interface{}) {
		if deathLetter, ok := evt.(*DeadLetterEvent); ok {
			if deathLetter.Sender != nil {
				logger.Error(deathLetter.Sender.GetID(), "DeathLetter Dest PID :%s Message:%+v", deathLetter.PID.String(), deathLetter.Message)
			} else {
				logger.Error(0, "DeathLetter Dest PID: %s Message:%+v", deathLetter.PID.String(), deathLetter.Message)
			}
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

//DeadLetterEvent desc
//@struct DeadLetterEvent desc:  Death news
type DeadLetterEvent struct {
	PID     *PID
	Message interface{}
	Sender  *PID
}

//SendUsrMessage desc
//@method SendUsrMessage desc: send a user message to death subscribe
//@param (PID) dest actor ID
//@param (interface{}) message
func (*deathLetterProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	})
}

//SendSysMessage desc
//@method SendSysMessage desc: send a system message to death subscribe
//@param (PID) dest actor ID
//@param (interface{}) message
func (*deathLetterProcess) SendSysMessage(pid *PID, message interface{}) {
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: message,
	})
}

//OverloadUsrMessage desc
//@method OverloadUsrMessage desc: user mesage queue overload
//@return (int) user mesage queue overload of number
func (*deathLetterProcess) OverloadUsrMessage() int {
	return 0
}

//Stop desc
//@method Stop desc: send stop message
//@param (PID) dest actor ID
func (slf *deathLetterProcess) Stop(pid *PID) {
	slf.SendSysMessage(pid, stopMessage)
}
