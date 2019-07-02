package bootstrap

import (
  "magicNet/engine/monitor"
  "magicNet/engine/preset/preset_hook"
  "magicNet/engine"
)

type Lanuch struct {
  inst engine.Framework
}

func NewLanuch(fk engine.Framework) *Lanuch {
  return &Lanuch{fk}
}

func (lch *Lanuch) Do() {
  engine.SetEngineInitHook(&preset_hook.DefaultEngineHook{})
  monitor.SetMonitorInitHook(&preset_hook.DefaultMonitorHook{})

  if lch.inst.Start() != 0 {
		goto lable_shutdown
	}

  lch.inst.Loop()

lable_shutdown:
  lch.inst.Shutdown()
}
