package actor

import (
  "magicNet/engine/util"
)

/***************************************
* 高15位表示服务器地址 | 低17表示PID编号 *
****************************************/
const (
  pidMask = 0x1ffff
  pidMax  = pidMask
)

type PID struct {
  address uint32
}

func (pid *PID)Compare(opid *PID) bool {
  if pid.address == opid.address {
    return true
  }
  return false
}

func (pid *PID)Hash(v uint32) uint32 {
  util.Assert(util.IsPower(int(v)), "hash miscalculation must be a power of 2")
  return pid.address & (v - 1)
}
