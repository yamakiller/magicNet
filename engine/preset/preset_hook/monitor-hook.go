package preset_hook

import "magicNet/engine/hook"
import "magicNet/engine/preset/preset_function"

type DefaultMonitorHook struct {
}

func (dmhook *DefaultMonitorHook)Initialize() bool {
  preset_function.InitializeAuth2()
  preset_function.RegisterMonitorBusiness()
  return true
}

func (dmhook *DefaultMonitorHook)Finalize() {

}

var (
	_ hook.InitializeHook = &DefaultMonitorHook{}
)
