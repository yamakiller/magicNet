package core

import (
	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicLibs/files"
)

//DefaultBoot deac
//@Struct DefaultBoost
//@Member logger
type DefaultBoot struct {
}

//Initial doc
//@Summary Initialization system
//@Method Initial
//@Return error Initialization fail returns error
func (slf *DefaultBoot) Initial() error {

	//read coroutine pool config
	coEnvPath := args.Instance().GetString("-c", "./config/coroutine_pool.json")
	coDeploy := coroutine.NewDefault()
	envs.Instance().Load(coroutine.EnvKey, coEnvPath, coDeploy)

	//read project root directed
	rootDir := args.Instance().GetString("root", "")

	//startup coroutine pool
	coroutine.Instance().Start(coDeploy.Max, coDeploy.Min, coDeploy.Task)

	if rootDir != "" {
		files.Instance().WithRoot(rootDir)
	}

	return nil
}

//Destory doc
//@Summary destory system reouse
//@Method Destory
func (slf *DefaultBoot) Destory() {
}
