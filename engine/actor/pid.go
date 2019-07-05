package actor

import (
  "sync/atomic"
  "unsafe"
)

/***************************************
* 高15位表示服务器地址 | 低17表示PID编号 *
****************************************/
const (
  pidMask    = 0x1ffff
  pidMax     = pidMask
  pidKeyBit  = 17
)

type PID struct {
  id uint32
  p  *Process
}

func (pid *PID) Address() uint32 {
  return pid.id >> pidKeyBit
}

func (pid *PID) Key() uint32 {
  return pid.id & pidMask
}

func (pid *PID) Compare(opid *PID) bool {
  if pid.id == opid.id {
    return true
  }
  return false
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

func NewPID() *PID {
  pid := &PID{}; GlobalRegistry.Register(pid)
  return pid
}

func pidFromId(id string, p *PID) {
  /*if ((id >> pidKeyBit) == GlobalRegistry.GetLocalAddress()) {

  }*/
}

func pidIsRemote(id uint32) bool {
  if  ((id >> pidKeyBit) == GlobalRegistry.GetLocalAddress()) {
    return false
  } else {
    return true
  }
}
