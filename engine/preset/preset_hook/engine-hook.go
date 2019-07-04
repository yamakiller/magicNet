package preset_hook

import(
  "magicNet/engine/hook"
)

type DefaultEngineHook struct {
}

func (dlhook *DefaultEngineHook)Initialize() bool {
  return true
}

func (dlhook *DefaultEngineHook)Finalize() {
}

var (
	_ hook.InitializeHook = &DefaultEngineHook{}
)
