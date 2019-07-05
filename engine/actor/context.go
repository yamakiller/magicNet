package actor

import "time"

type infoPart interface {
  Self() *PID

  Actor() Actor
}

type Context interface {
  infoPart
  basePart
  messagePart
  senderPart
  receiverPart
  stopperPart
}

type SenderContext interface {
  infoPart
  senderPart
  messagePart
}

type ReceiverContext interface {
  infoPart
  receiverPart
  messagePart
}

type basePart interface {
  ReceiveTimeout() time.Duration

  Respond(response interface{})

  //将当前的消息，存放到stack上
  Stash()

  //注册监视器
  Watch(pid *PID)

  //注销监视器
  Unwatch(pid *PID)

  //设置定时器
  SetReceiveTimeout(d time.Duration)

  //取消定时器
	CancelReceiveTimeout()

  //将当前消息转发给指定的PID
  Forward(pid *PID)
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
}

type receiverPart interface{
  Receive(pack *MessagePack)
}

type stopperPart interface {
  Stop(pid *PID)

  Poison(pid *PID)
}
