package test


import (
	"magicNet/engine"
	"magicNet/logger"
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
  lanuch := engine.NewLanuch(engine.Framework{}, &testHook{})
  lanuch.Do()
}

func TestLanuchHookDefault() {
	lanuch := engine.NewLanuch(engine.Framework{}, nil)
  lanuch.Do()
}
