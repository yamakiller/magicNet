package actor

import "github.com/yamakiller/magicNet/engine/mailbox"

// AutoReceiveMessage Receive messages automatically
type AutoReceiveMessage interface {
	AutoReceiveMessage()
}

// NotInfluenceReceiveTimeout ：不影响接收超时
type NotInfluenceReceiveTimeout interface {
	NotInfluenceReceiveTimeout()
}

// SystemMessage ：系统消息
type SystemMessage interface {
	SystemMessage()
}

type continuation struct {
	message interface{}
	f       func()
}

// AutoReceiveMessage  Stop message Defines the AutoReceiveMessage function object
func (*Stopping) AutoReceiveMessage() {}

// AutoReceiveMessage Stopped message Defines the AutoReceiveMessage function object
func (*Stopped) AutoReceiveMessage() {}

// SystemMessage Started message Defining the SystemMessage function object
func (*Started) SystemMessage() {}

// SystemMessage Stop message Define SystemMessage function object
func (*Stop) SystemMessage() {}

// SystemMessage Enter observation/join observation message Define SystemMessage function object
func (*Watch) SystemMessage() {}

// SystemMessage Unobserve/Export Watch Message Define SystemMessage Function Object
func (*Unwatch) SystemMessage() {}

// SystemMessage Terminate message Define SystemMessage function object
func (*Terminated) SystemMessage() {}

// SystemMessage Proxy call Defines the SystemMessage function object
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
