package actor

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月05日 16:35:59
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 17:21:40
 * @Description:  Context 接口
 */

import "time"

type infoPart interface {
	Self() *PID

	Actor() Actor
}

// Context : Context 基础接口
type Context interface {
	infoPart
	basePart
	messagePart
	senderPart
	receiverPart
	makerPart
	stopperPart
}

// SenderContext ： 发送者Context基础接口
type SenderContext interface {
	infoPart
	senderPart
	messagePart
}

// ReceiverContext : 接收者Context基础接口
type ReceiverContext interface {
	infoPart
	receiverPart
	messagePart
}

// MakerContext : 创建者Context基础接口
type MakerContext interface {
	infoPart
	makerPart
}

// basePart : 所有对象的基础接口
type basePart interface {
	//ReceiveTimeout() time.Duration

	Respond(response interface{})

	//将当前的消息，存放到stack上
	Stash()

	//注册监视器
	Watch(pid *PID)

	//注销监视器
	Unwatch(pid *PID)

	//设置定时器
	//SetReceiveTimeout(d time.Duration)

	//取消定时器
	//CancelReceiveTimeout()

	//将当前消息转发给指定的PID
	Forward(pid *PID)

	AwaitFuture(f *Future, continuation func(res interface{}, err error))
}

type messagePart interface {
	Message() interface{}

	MessageHeader() ReadOnlyMessageHeader
}

type senderPart interface {
	Sender() *PID

	Send(pid *PID, message interface{})

	Request(pid *PID, message interface{})

	RequestWithCustomSender(pid *PID, message interface{}, sender *PID)

	RequestFuture(pid *PID, message interface{}, timeout time.Duration) *Future
}

type receiverPart interface {
	Receive(pack *MessagePack)
}

type makerPart interface {
	Make(agnet *Agnets) *PID

	MakeNamed(agnet *Agnets, name string) (*PID, error)
}

type stopperPart interface {
	Stop(pid *PID)

	StopFuture(pid *PID) *Future

	Kill(pid *PID)

	KillFuture(pid *PID) *Future
}
