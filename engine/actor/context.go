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

// Context Context interface
type Context interface {
	infoPart
	basePart
	messagePart
	senderPart
	receiverPart
	makerPart
	stopperPart
}

// SenderContext Sender Context Base Interface
type SenderContext interface {
	infoPart
	senderPart
	messagePart
}

// ReceiverContext Receiver Context Basic Interface
type ReceiverContext interface {
	infoPart
	receiverPart
	messagePart
}

// MakerContext  Creator Context base interface
type MakerContext interface {
	infoPart
	makerPart
}

// basePart Basic interface for all objects
type basePart interface {
	//ReceiveTimeout() time.Duration

	Respond(response interface{})

	//Store the current message on the stack
	Stash()

	//Registration monitor
	Watch(pid *PID)

	//注销监视器
	Unwatch(pid *PID)

	//Set timer
	//SetReceiveTimeout(d time.Duration)

	//Cancel timer
	//CancelReceiveTimeout()

	//Forward the current message to the specified PID
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
