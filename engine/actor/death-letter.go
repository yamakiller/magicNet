package actor

import (
	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicNet/engine/evtchan"
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

//DeadLetterEvent doc
//@Struct DeadLetterEvent @Summary  Death news
type DeadLetterEvent struct {
	PID     *PID
	Message interface{}
	Sender  *PID
}

//SendUsrMessage doc
//@Method SendUsrMessage @Summary send a user message to death subscribe
//@Param (PID) dest actor ID
//@Param (interface{}) message
func (*deathLetterProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	})
}

//SendSysMessage doc
//@Method SendSysMessage @Summary send a system message to death subscribe
//@Param (PID) dest actor ID
//@Param (interface{}) message
func (*deathLetterProcess) SendSysMessage(pid *PID, message interface{}) {
	evtchan.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: message,
	})
}

//OverloadUsrMessage doc
//@Method OverloadUsrMessage @Summary user mesage queue overload
//@Return (int) user mesage queue overload of number
func (*deathLetterProcess) OverloadUsrMessage() int {
	return 0
}

//Stop doc
//@Method Stop @Summary send stop message
//@Param (PID) dest actor ID
func (slf *deathLetterProcess) Stop(pid *PID) {
	slf.SendSysMessage(pid, stopMessage)
}
