package launch

import (
	"magicNet/core/debug"
	"magicNet/core/frame"
)

// Launch : 系统启动器
func Launch(f frame.MakeFrame) {
	defer debug.Trace()
	fme := f()
	if err := fme.Init(); err != nil {
		panic(err)
	}

	if err := fme.LoadEnv(); err != nil {
		panic(err)
	}

	if err := fme.InitService(); err != nil {
		panic(err)
	}

	for {
		if fme.Wait() == -1 {
			break
		}
	}

	fme.CloseService()
	fme.UnLoadEnv()
	fme.Destory()
}
