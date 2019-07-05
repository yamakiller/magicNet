package actor


type Actor interface {
  Receive(c Context)
}

type ActorFunc func(c Context)

func(f ActorFunc) Receive(c Context) {
  f(c)
}

type ReceiverFunc func(c ReceiverContext, pack *MessagePack)

type SenderFunc func(c SenderContext, pack *MessagePack)

type ContextDecoratorFunc func(ctx Context) Context
