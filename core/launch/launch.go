package launch

import (
	"fmt"
	"os"

	"github.com/yamakiller/magicNet/core/debug"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/timer"
)

//Launch desc
//@method Launch desc: Start function
//@param (frame.MakeFrame) Start framework
func Launch(f frame.MakeFrame) {
	fme := f()
	fme.Option()

	debugTrace := debug.TraceDebug{}
	debugTrace.Start()
	defer debugTrace.Stop()

	timer.StartService()

	if err := fme.LoadEnv(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err := fme.Init(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err := fme.InitService(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fme.EnterLoop()

	for {
		if fme.Wait() == -1 {
			break
		}
	}

	fme.CloseService()
	fme.Destory()
	fme.UnLoadEnv()
	timer.StopService()
}
