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
  spawnerPart
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

type SpawnerContext interface {
  infoPart
  spawnerPart
}

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

type receiverPart interface{
  Receive(pack *MessagePack)
}

type spawnerPart interface {

	//Spawn(props *Props) *PID

	//SpawnPrefix(props *Props, prefix string) *PID

	//SpawnNamed(props *Props, id string) (*PID, error)
}

type stopperPart interface {
  Stop(pid *PID)

  StopFuture(pid *PID) *Future

  Kill(pid *PID)

  KillFuture(pid *PID) *Future
}
