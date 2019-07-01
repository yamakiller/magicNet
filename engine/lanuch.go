package engine

import "magicNet/normal"

type Lanuch struct {
  inst Framework
  hook LanuchHook
}

type LanuchHook interface {
  Initialize() bool
  Finalize()
}

func NewLanuch(fk Framework, hk LanuchHook) *Lanuch {
  if hk == nil {
    hk = &normal.DefaultLanuchHook{}
  }
  return &Lanuch{fk, hk}
}

func (lch *Lanuch) Do() {
  if lch.inst.start() != 0 {
		goto lable_shutdown
	}

  if lch.hook != nil &&
  !lch.hook.Initialize() {
    goto lable_finalize
  }

  lch.inst.loop()

lable_finalize:
  if lch.hook != nil  {
    lch.hook.Finalize()
  }
lable_shutdown:
  lch.inst.shutdown()
}
