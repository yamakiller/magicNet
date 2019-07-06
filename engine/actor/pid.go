package actor

import (
  "sync/atomic"
  "unsafe"
  "strings"
  "time"
  "magicNet/engine/util"
)

/***************************************
* 高15位表示服务器地址 | 低17表示PID编号 *
****************************************/
const (
  pidMask    = 0x1ffff
  pidMax     = pidMask
  pidKeyBit  = 17
)

func idToHex(u uint32) string {
  const (
    digits = "0123456789ABCDEF"
  )

  var str[10]byte
  str[0] = '$'
  var i uint32
  for i = 0;i < 8; i++ {
    str[i + 1] = digits[(u >> ((7 - i) * 4)) & 0xf]
  }

  return string(str[:8])
}

func pidFromId(id string, p *PID) {
  var i uint32
  var addr uint32 = 0
  var len = uint32(strings.Count(id, "") - 1)
  for i = 1;i < len;i++ {
    c := id[i]
    if c >= '0' && c <= '9' {
      c = c - '0'
    } else if c >= 'a' && c <= 'f' {
      c = c - 'a' + 10
    } else if (c >= 'A' && c <= 'F') {
      c = c - 'A' + 10
    } else {
      util.Assert(false, "Id unknown character")
    }
    addr = addr * 16 + uint32(c)
  }
}

func pidIsRemote(id uint32) bool {
  if  ((id >> pidKeyBit) == GlobalRegistry.GetLocalAddress()) {
    return false
  } else {
    return true
  }
}

type PID struct {
  Id uint32
  p  *Process
}

func (pid *PID) Address() uint32 {
  return pid.Id >> pidKeyBit
}

func (pid *PID) Key() uint32 {
  return pid.Id & pidMask
}

func (pid *PID) ref() Process {
  p := (*Process)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p))))
  if p != nil  {
    if l, ok := (*p).(*ActorProcess); ok && atomic.LoadInt32(&l.death) == 1 {
      atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), nil)
    } else {
      return *p
    }
  }

  ref, exits :=  GlobalRegistry.Get(pid)
  if exits {
     atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), unsafe.Pointer(&ref))
  }

  return ref
}

func (pid *PID) sendUsrMessage(message interface{}) {
  pid.ref().SendUsrMessage(pid, message)
}

func (pid *PID) sendSysMessage(message interface{}) {
  pid.ref().SendSysMessage(pid, message)
}

func (pid *PID) String() string {
  return ""
}

func (pid *PID) Stop() {
  pid.ref().Stop(pid)
}

func NewPID() *PID {
  pid := &PID{}; GlobalRegistry.Register(pid)
  return pid
}

func (pid *PID) Tell(message interface{}) {
  ctx := DefaultForwardContext
  ctx.Send(pid, message)
}

func (pid *PID) Request(message interface{}, responseTo *PID) {
  ctx := DefaultForwardContext
  ctx.RequestWithCustomSender(pid, message, responseTo)
}

func (pid *PID) RequestFuture(message interface{}, timeOut time.Duration) *Future {
  ctx := DefaultForwardContext
  return ctx.RequestFuture(pid, message, timeOut)
}

func (pid *PID) StopFuture() *Future {
  future := NewFuture(10 * time.Second)

  pid.sendSysMessage(&Watch{Watcher: future.pid})
  pid.Stop()
  return future
}

func (pid *PID) StopWait() {
  pid.StopFuture().Wait()
}
