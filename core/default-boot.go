package core

import (
	"runtime"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicLibs/files"
	"github.com/yamakiller/magicLibs/logger"
)

//DefaultBoot deac
//@strudt DefaultBoost
//@Member (logger)
type DefaultBoot struct {
	_log logger.Logger
}

//Initial desc
//@Method Initial desc: Initialization system
//@Return (error) Initialization fail returns error
func (slf *DefaultBoot) Initial() error {
	//read log module config
	logEnvPath := args.Instance().GetString("-l", "./config/log.json")
	logDeploy := logger.NewDefault()
	envs.Instance().Load(logger.EnvKey, logEnvPath, logDeploy)

	//read coroutine pool config
	coEnvPath := args.Instance().GetString("-c", "./config/coroutine_pool.json")
	coDeploy := coroutine.NewDefault()
	envs.Instance().Load(coroutine.EnvKey, coEnvPath, coDeploy)

	//read project root directed
	rootDir := args.Instance().GetString("root", "")

	//startup coroutine pool
	coroutine.Instance().Start(coDeploy.Max, coDeploy.Min, coDeploy.Task)
	//startup logger
	slf._log = logger.New(func() logger.Logger {
		l := logger.LogContext{}
		l.SetFilPath(logDeploy.LogPath)
		l.SetHandle(logrus.New())
		l.SetMailMax(logDeploy.LogSize)
		l.SetLevel(logrus.Level(logDeploy.LogLevel))

		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		if runtime.GOOS == "windows" {
			formatter.DisableColors = true
		} else {
			formatter.SetColorScheme(&prefixed.ColorScheme{
				PrefixStyle: "blue+b"})
		}
		l.SetFormatter(formatter)
		l.Initial()
		l.Redirect()
		return &l
	})

	logger.WithDefault(slf._log)
	slf._log.Mount()

	if rootDir != "" {
		files.Instance().WithRoot(rootDir)
	}

	return nil
}

//Destory desc
//@Method Destory desc: destory system reouse
func (slf *DefaultBoot) Destory() {
	logger.Info(0, "Destory")
	if slf._log != nil {
		slf._log.Close()
		slf._log = nil
	}
}
