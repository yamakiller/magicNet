package actor

import "github.com/yamakiller/magicNet/engine/mailbox"

// AutoReceiveMessage ： 自动接收消息
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

// AutoReceiveMessage  : 停止中消息 定义AutoReceiveMessage 函数对象
func (*Stopping) AutoReceiveMessage() {}

// AutoReceiveMessage : 已停止消息  定义AutoReceiveMessage 函数对象
func (*Stopped) AutoReceiveMessage() {}

// SystemMessage : 已经开始消息   定义SystemMessage 函数对象
func (*Started) SystemMessage() {}

// SystemMessage : 停止消息   定义SystemMessage 函数对象
func (*Stop) SystemMessage() {}

// SystemMessage : 进入观察/加入观察消息   定义SystemMessage 函数对象
func (*Watch) SystemMessage() {}

// SystemMessage : 取消观察/推出观察消息   定义SystemMessage 函数对象
func (*Unwatch) SystemMessage() {}

// SystemMessage : 终止消息   定义SystemMessage 函数对象
func (*Terminated) SystemMessage() {}

// SystemMessage : 代理调用   定义SystemMessage 函数对象
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
