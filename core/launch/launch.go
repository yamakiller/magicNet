package launch

import (
	"flag"
	"fmt"
	"os"

	"github.com/yamakiller/magicNet/core/debug"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/timer"
)

// Launch : 系统启动器
func Launch(f frame.MakeFrame) {
	fme := f()
	fme.VarValue()
	flag.Parse()
	fme.LineOption()

	debugTrace := debug.TraceDebug{}
	debugTrace.Start()
	defer debugTrace.Stop()

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
