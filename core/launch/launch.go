package launch

import (
	"magicNet/core/debug"
	"magicNet/core/frame"
)

// Launch : 系统启动器
func Launch(f frame.MakeFrame) {
	defer debug.Trace()
	fme := f()
	if !fme.Start() {
		goto l_start_lable
	}

	if !fme.LoadEnv() {
		goto l_env_label
	}

	if !fme.InitService() {
		goto l_srv_label
	}

	for {
		if fme.Wait() == -1 {
			break
		}
	}

l_srv_label:
	fme.CloseService()
l_env_label:
	fme.UnLoadEnv()
l_start_lable:
	fme.Shutdown()
}
