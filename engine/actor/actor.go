package actor

import (
  "sync/atomic"
)

type IActor interface {
  Pid() *PID
  Dispone()
  Stop()
}

type Actor struct {
  pid PID
  dead int32
}

func (a *Actor)Pid() *PID {
  return &a.pid
}

func (a *Actor)Dispone() {

}

func (a *Actor)Stop() {
  atomic.StoreInt32(&a.dead, 1)
}

var (
  _ IActor = &Actor{}
)
