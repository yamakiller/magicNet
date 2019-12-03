package boot

import (
	"fmt"

	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicNet/core/debug"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/timer"
)

//Launch doc
//@Method Launch @Summary Start function
//@Param (frame.SpawnFrame) Start framework
func Launch(f frame.SpawnFrame) {
	args.Instance().Parse()
	dTrace := debug.TraceDebug{}
	dTrace.Start()
	timer.StartService()
	defer func() {
		dTrace.Stop()
		timer.StopService()
	}()

	fme := f()
	if err := fme.Initial(); err != nil {
		fmt.Println(err)
		goto end
	}

	if err := fme.InitService(); err != nil {
		fmt.Println(err)
		goto end
	}

	fme.Enter()

	for {
		if fme.Wait() == -1 {
			break
		}
	}
end:
	fme.CloseService()
	fme.Destory()
}
