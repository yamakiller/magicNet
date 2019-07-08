package actor

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月02日 16:43:48
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 16:37:34
 * @Description: Actor 基础接口
 */

// Actor : 基础接口
type Actor interface {
	Receive(c Context)
}

// MakeActor : 创建 Actor函数
type MakeActor func() Actor

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
