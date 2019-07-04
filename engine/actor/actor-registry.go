package actor

import (
  "sync"
  "sync/atomic"
  "runtime"
  "magicNet/engine/util"
)


type ActorRegistry struct {
  localAddress uint32
  localSequence uint32
  localItem []IActor
  localItemMutex sync.RWMutex
}

func (a *ActorRegistry)SetLocalAddress(addr uint32) {
  a.localAddress = addr
}

func (a* ActorRegistry)Register(actor IActor) PID {
  a.localItemMutex.Lock()
  for {
    var i uint32
    currentNum := uint32(len(a.localItem))
    for i = 0;i < currentNum;i++ {
       address := ((i + a.localSequence) & pidMask)
       hash := address & (currentNum - 1)
       if a.localItem[hash] == nil {
          a.localItem[hash] = actor
          a.localSequence = address + 1
          a.localItemMutex.Unlock()

          return PID{address & a.localAddress}
       }
    }

    newnum := (currentNum << 1)
    util.Assert(newnum <= pidMax, "actor number overflow")
    newItem := make([]IActor, newnum)

    for i=0;i < currentNum;i++ {
      if newItem[i] == nil {
        continue
      }

      hash := newItem[i].Pid().Hash(newnum)
      if (hash == i) {
        continue
      }
      newItem[hash] = a.localItem[i]
    }

    a.localItem = newItem
  }
}

func (a *ActorRegistry)UnRegister(pid *PID) bool {
  a.localItemMutex.Lock()
  defer a.localItemMutex.Unlock()
  hash := pid.address & uint32(len(a.localItem) - 1)
  if a.localItem[hash] != nil && a.localItem[hash].Pid().Compare(pid)  {
    ref := a.localItem[hash]
    if l, ok := ref.(*Actor); ok {
      atomic.StoreInt32(&l.dead, 1)
    }
    runtime.SetFinalizer(&ref, func(a IActor){a.Dispone()})
    a.localItem[hash] = nil
    return true
  }
  return false
}

func (a *ActorRegistry)Get(pid *PID) IActor {
  if pid == nil {
    return nil
  }

  return nil
}
