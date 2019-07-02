package preset_hook

import(
  //"magicNet/logger"
  //"magicNet/engine/preset/preset_function"
)

type DefaultEngineHook struct {
}

func (dlhook *DefaultEngineHook)Initialize() bool {
  return true
}

func (dlhook *DefaultEngineHook)Finalize() {
}
