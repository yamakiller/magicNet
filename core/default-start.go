package core

import (
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/yamakiller/magicLibs/files"

	"github.com/yamakiller/magicLibs/args"

	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicLibs/logger"
)

//DefaultStart desc
//@struct DefaultStart desc: Default launcher
type DefaultStart struct {
	sysLogger logger.Logger
}

//Init desc
//@method Init desc: Initializing system
//@return (error) Initializing fail Returns error
func (slf *DefaultStart) Init() error {

	coDeplay := envs.Instance().Get(coroutine.EnvKey).(*coroutine.Deploy)
	logDeplay := envs.Instance().Get(logger.EnvKey).(*logger.LogDeploy)
	rootDir := args.Instance().GetString("root", "")

	coroutine.Instance().Start(coDeplay.Max, coDeplay.Min, coDeplay.Task)

	slf.sysLogger = logger.New(func() logger.Logger {
		l := logger.LogContext{}
		l.SetFilPath(logDeplay.LogPath)
		l.SetHandle(logrus.New())
		l.SetMailMax(logDeplay.LogSize)
		l.SetLevel(logrus.Level(logDeplay.LogLevel))

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

	logger.WithDefault(slf.sysLogger)
	slf.sysLogger.Mount()

	if rootDir != "" {
		files.Instance().WithRoot(rootDir)
	}

	return nil
}

//Destory desc
//@method Destory desc: Destruction processing
func (slf *DefaultStart) Destory() {
	if slf.sysLogger != nil {
		slf.sysLogger.Close()
	}

	files.Instance().Close()
	coroutine.Instance().Stop()
}
