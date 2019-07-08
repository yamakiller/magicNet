package testing

import (
	"fmt"
	"magicNet/bootstrap"
	"magicNet/engine"
	"magicNet/engine/hook"
	"magicNet/engine/logger"
	"magicNet/engine/monitor"
)

type testHook struct {
}

func (t *testHook) Initialize() bool {
	logger.Info(0, "test hook Initialize")
	return false
}

func (t *testHook) Finalize() {
	logger.Info(0, "test hook Finalize")
}

var (
	_ hook.InitializeHook = &testHook{}
)

// TestEmpty : xxxxx
func TestEmpty() {

}

// TestWait :
func TestWait() {
	fmt.Println("已经执行到末尾")
	var ispass bool
	fmt.Scanln(&ispass)
	//fmt.Println("程序结束\n")
}

// TestLanuchHook : xxxx
func TestLanuchHook() {
	engine.SetEngineInitHook(&testHook{})
	monitor.SetMonitorInitHook(&testHook{})
	lanuch := bootstrap.NewLanuch(engine.Framework{})
	lanuch.Do()
}

// TestLanuchHookDefault : xxxxxxx
func TestLanuchHookDefault() {
	lanuch := bootstrap.NewLanuch(engine.Framework{})
	lanuch.Do()
}
