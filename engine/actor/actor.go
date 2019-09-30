package actor

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月02日 16:43:48
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 16:37:34
 * @Description: Actor 基础接口
 */

// Actor Basic interface
type Actor interface {
	Receive(c Context)
}

// MakeActor Create an Actor function
type MakeActor func() Actor

// AtrFunc Actor Receive Agent Function
type AtrFunc func(c Context)

// Receive : Actor receive function shell
func (f AtrFunc) Receive(c Context) {
	f(c)
}

// ReceiverFunc : Actor receive function definition
type ReceiverFunc func(c ReceiverContext, pack *MessagePack)

// SenderFunc : Actor sender function definition
type SenderFunc func(c SenderContext, target *PID, pack *MessagePack)
