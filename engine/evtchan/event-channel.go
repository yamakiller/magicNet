package evtchan

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月06日 10:12:19
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 15:31:45
 * @Description:事件通道 [同步]
 */

import "sync"

// Predicate ： 先决条件函数
type Predicate func(evt interface{}) bool

var ec = &EventChannel{}

// Subscribe ：默认通道：订阅
func Subscribe(fn func(evt interface{})) *Subscription {
	return ec.Subscribe(fn)
}

// UnSubscribe ：默认通道：取消订阅
func UnSubscribe(sub *Subscription) {
	ec.UnSubscribe(sub)
}

// Publish ：发布默认事件
func Publish(evt interface{}) {
	ec.Publish(evt)
}

// EventChannel ：事件通道
type EventChannel struct {
	sync.RWMutex
	subscriptions []*Subscription
}

// Subscribe ：订阅
func (ec *EventChannel) Subscribe(fn func(evt interface{})) *Subscription {
	ec.Lock()
	defer ec.Unlock()

	sub := &Subscription{
		ec: ec,
		i:  len(ec.subscriptions),
		fn: fn,
	}

	ec.subscriptions = append(ec.subscriptions, sub)
	return sub
}

// UnSubscribe ： 取消订阅
func (ec *EventChannel) UnSubscribe(sub *Subscription) {
	if sub.i == -1 {
		return
	}

	ec.Lock()
	defer ec.Unlock()

	i := sub.i
	l := len(ec.subscriptions) - 1
	ec.subscriptions[i] = ec.subscriptions[l]
	ec.subscriptions[i].i = i
	ec.subscriptions[l] = nil
	ec.subscriptions = ec.subscriptions[:l]
	sub.i = -1

	if len(ec.subscriptions) == 0 {
		ec.subscriptions = nil
	}
}

// Publish ：发布
func (ec *EventChannel) Publish(evt interface{}) {
	ec.RLock()
	defer ec.RUnlock()

	for _, s := range ec.subscriptions {
		if s.p == nil || s.p(evt) {
			s.fn(evt)
		}
	}
}

// Subscription ：订阅器
type Subscription struct {
	ec *EventChannel
	i  int
	fn func(event interface{})
	p  Predicate
}

// WithPredicate ：设置先决条件函数
func (s *Subscription) WithPredicate(p Predicate) *Subscription {
	s.ec.Lock()
	s.p = p
	s.ec.Unlock()
	return s
}
