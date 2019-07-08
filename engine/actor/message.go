package actor

import "magicNet/engine/mailbox"

type AutoReceiveMessage interface {
	AutoReceiveMessage()
}

type NotInfluenceReceiveTimeout interface {
	NotInfluenceReceiveTimeout()
}

type SystemMessage interface {
	SystemMessage()
}

type continuation struct {
	message interface{}
	f       func()
}

func (*Stopping) AutoReceiveMessage() {}
func (*Stopped) AutoReceiveMessage()  {}

func (*Started) SystemMessage()      {}
func (*Stop) SystemMessage()         {}
func (*Watch) SystemMessage()        {}
func (*Unwatch) SystemMessage()      {}
func (*Terminated) SystemMessage()   {}
func (*continuation) SystemMessage() {}

var (
	stoppingMessage       interface{} = &Stopping{}
	stoppedMessage        interface{} = &Stopped{}
	receiveTimeoutMessage interface{} = &ReceiveTimeout{}
)

var (
	startedMessage        interface{} = &Started{}
	stopMessage           interface{} = &Stop{}
	resumeMailboxMessage  interface{} = &mailbox.ResumeMailbox{}
	suspendMailboxMessage interface{} = &mailbox.SuspendMailbox{}
)
