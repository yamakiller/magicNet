package test


import (
	"magicNet/engine"
	"magicNet/engine/monitor"
	"magicNet/engine/logger"
	"magicNet/bootstrap"
)

type testHook struct {

}

func (t *testHook)Initialize() bool {
	logger.Info(0, "test hook Initialize")
	return false
}

func (t *testHook)Finalize() {
  	logger.Info(0, "test hook Finalize")
}

func TestEmpty() {

}

func TestLanuchHook() {
	engine.SetEngineInitHook(&testHook{})
	monitor.SetMonitorInitHook(&testHook{})
	lanuch := bootstrap.NewLanuch(engine.Framework{})
  lanuch.Do()
}

func TestLanuchHookDefault() {
	lanuch := bootstrap.NewLanuch(engine.Framework{})
  lanuch.Do()
}
