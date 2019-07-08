package actor

// Actor : 基础接口
type Actor interface {
	Receive(c Context)
}

// NewActor : 创建 Actor函数
type NewActor func() Actor

// ActorFunc : Actor 接收代理函数
type ActorFunc func(c Context)

// Receive : Actor 接收函数外壳
func (f ActorFunc) Receive(c Context) {
	f(c)
}

// ReceiverFunc : Actor 接收函数定义
type ReceiverFunc func(c ReceiverContext, pack *MessagePack)

// SenderFunc : Actor 发送者函数定义
type SenderFunc func(c SenderContext, target *PID, pack *MessagePack)

//type ContextDecoratorFunc func(ctx Context) Context
