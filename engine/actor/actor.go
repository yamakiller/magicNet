package actor

// Actor : 基础接口
type Actor interface {
	Receive(c Context)
}

// NewActor : 创建 Actor函数
type NewActor func() Actor

// AtrFunc : Actor 接收代理函数
type AtrFunc func(c Context)

// Receive : Actor 接收函数外壳
func (f AtrFunc) Receive(c Context) {
	f(c)
}

// ReceiverFunc : Actor 接收函数定义
type ReceiverFunc func(c ReceiverContext, pack *MessagePack)

// SenderFunc : Actor 发送者函数定义
type SenderFunc func(c SenderContext, target *PID, pack *MessagePack)

//type ContextDecoratorFunc func(ctx Context) Context
