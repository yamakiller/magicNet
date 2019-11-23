package core

import (
	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicLibs/logger"
)

//DefaultEnv desc
//@struct DefaultEnv desc: Default environment variable manager
type DefaultEnv struct {
}

//LoadEnv desc
//@method LoadEnv desc: Loading environment variables
func (slf *DefaultEnv) LoadEnv() error {

	logEnvPath := args.Instance().GetString("-l", "./config/log.json")
	logDeploy := logger.NewDefault()
	envs.Instance().Load(logger.EnvKey, logEnvPath, logDeploy)

	coEnvPath := args.Instance().GetString("-c", "./config/coroutine_pool.json")
	coDeploy := coroutine.NewDefault()
	envs.Instance().Load(coroutine.EnvKey, coEnvPath, coDeploy)

	return nil
}

//UnLoadEnv desc
//@method UnLoadEnv desc: Unload environment variable information
func (slf *DefaultEnv) UnLoadEnv() {
	envs.Instance().UnLoad()
}
