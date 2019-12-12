package core

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicLibs/files"
	"github.com/yamakiller/magicLibs/logger"
)

//DefaultBoot deac
//@Struct DefaultBoost
//@Member logger
type DefaultBoot struct {
	_log logger.Logger
}

//Initial doc
//@Summary Initialization system
//@Method Initial
//@Return error Initialization fail returns error
func (slf *DefaultBoot) Initial() error {
	//read log module config
	logPath := args.Instance().GetString("-log", "")
	logSize := args.Instance().GetInt("-logSize", 128)
	logLevel := logger.DEBUGLEVEL
	release := args.Instance().GetBoolean("-release", false)
	if release {
		logLevel = logger.INFOLEVEL
	}

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
		l.SetFilPath(logPath)
		l.SetHandle(logrus.New())
		l.SetMailMax(logSize)
		l.SetLevel(logrus.Level(logLevel))

		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		formatter.TimestampFormat = "2006-01-02 15:04:05"
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle:    "white+h",
			TimestampStyle: "black+h"})
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

//Destory doc
//@Summary destory system reouse
//@Method Destory
func (slf *DefaultBoot) Destory() {
	if slf._log != nil {
		slf._log.Close()
		slf._log = nil
	}
}
