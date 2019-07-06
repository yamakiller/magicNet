package eventchannel

import "sync"

type Predicate func(evt interface{}) bool

var ec = &EventChannel{}

func Subscribe(fn func(evt interface{})) *Subscription {
  return ec.Subscribe(fn)
}

func UnSubscribe(sub *Subscription) {
  ec.UnSubscribe(sub)
}

func Publish(evt interface{}) {
  ec.Publish(evt)
}

type EventChannel struct {
  sync.RWMutex
  subscriptions []* Subscription
}

func (ec *EventChannel) Subscribe(fn func(evt interface{})) *Subscription {
  ec.Lock()
  defer ec.Unlock()

  sub := &Subscription {
    ec: ec,
    i : len(ec.subscriptions),
    fn: fn,
  }

  ec.subscriptions = append(ec.subscriptions, sub)
  return sub
}

func (ec *EventChannel) UnSubscribe(sub *Subscription) {
  if sub.i == -1 {
    return
  }

  ec.Lock()
  defer ec.Unlock()

  i := sub.i
  l := len(ec.subscriptions) - 1
  ec.subscriptions[i]   = ec.subscriptions[l]
  ec.subscriptions[i].i = i
  ec.subscriptions[l]   = nil
  ec.subscriptions      = ec.subscriptions[:l]
  sub.i = -1

  if len(ec.subscriptions) == 0 {
    ec.subscriptions = nil
  }
}

func (ec *EventChannel) Publish(evt interface{}) {
  ec.RLock()
  defer ec.RUnlock()

  for _, s := range ec.subscriptions {
    if s.p == nil || s.p(evt) {
      s.fn(evt)
    }
  }
}

type Subscription struct {
  ec *EventChannel
  i   int
  fn  func(event interface{})
  p   Predicate
}

func (s *Subscription) WithPredicate(p Predicate) *Subscription {
	s.ec.Lock()
	s.p = p
	s.ec.Unlock()
	return s
}
