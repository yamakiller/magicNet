package normal

import(
  "magicNet/logger"
)

type DefaultLanuchHook struct {

}

func (dlhook *DefaultLanuchHook)Initialize() bool {
  initializeAuth2()
  return true
}

func (dlhook *DefaultLanuchHook)Finalize() {
  logger.Info(0, "default lanuch hook Finalize")
}
