package launch

import (
	"fmt"
	"magicNet/core/debug"
	"magicNet/core/frame"
	"magicNet/timer"
	"os"
)

// Launch : 系统启动器
func Launch(f frame.MakeFrame) {
	debugTrace := debug.TraceDebug{}
	debugTrace.Start()
	defer debugTrace.Stop()
	fme := f()
	timer.StartService()

	if err := fme.Init(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err := fme.LoadEnv(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if err := fme.InitService(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for {
		if fme.Wait() == -1 {
			break
		}
	}

	fme.CloseService()
	fme.UnLoadEnv()
	fme.Destory()
	timer.StopService()
}
