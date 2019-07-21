package main

import (
	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/core/launch"
)

func main() {
	launch.Launch(func() frame.Framework {
		return &core.DefaultFrame{}
	})
	//testing.TestEmpty()
	//testing.TestTimer()
	//testing.TestChan()
	//testing.TestRectPoint()
	//testing.TestJSStack()
	//testing.TestDir()
	//testing.TestActorContext()
	//test.TestLanuchHook()
	//test.TestLanuchHookDefault()
	//testing.TestWait()
}
