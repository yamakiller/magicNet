package actor

type deathLetterProcess struct {}

var (
  deathLetter Process = &deathLetterProcess{}
)

func init() {

}

type DeadLetterEvent struct {
	PID     *PID        // The invalid process, to which the message was sent
	Message interface{} // The message that could not be delivered
	Sender  *PID        // the process that sent the Message
}

func (*deathLetterProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	/*eventstream.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	})*/
}

func (*deathLetterProcess) SendSysMessage(pid *PID, message interface{}) {
	/*eventstream.Publish(&DeadLetterEvent{
		PID:     pid,
		Message: message,
	})*/
}

func (ref *deathLetterProcess) Stop(pid *PID) {
	ref.SendSysMessage(pid, stopMessage)
}
