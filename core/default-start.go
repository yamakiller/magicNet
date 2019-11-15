package core

import (
	"runtime"

	"github.com/yamakiller/magicLibs/files"

	"github.com/yamakiller/magicLibs/args"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/coroutine"
	"github.com/yamakiller/magicLibs/envs"
	"github.com/yamakiller/magicNet/engine/logger"
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

	coDeplay := envs.Instance().Get("coroutine pool").(*coroutine.Deploy)
	logDeplay := envs.Instance().Get("log").(*logger.LogDeploy)
	rootDir := args.Instance().GetString("root", "")

	coroutine.Instance().Start(coDeplay.Max, coDeplay.Min, coDeplay.Task)

	slf.sysLogger = logger.New(func() logger.Logger {
		l := logger.LogContext{FilName: logDeplay.LogPath,
			LogHandle:  logrus.New(),
			LogMailbox: make(chan logger.Event, logDeplay.LogSize),
			LogStop:    make(chan struct{})}

		l.LogHandle.SetLevel(logrus.Level(logDeplay.LogLevel))

		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		if runtime.GOOS == "windows" {
			formatter.DisableColors = true
		} else {
			formatter.SetColorScheme(&prefixed.ColorScheme{
				PrefixStyle: "blue+b"})
		}
		l.LogHandle.SetFormatter(formatter)
		l.Redirect()
		return &l
	})

	coroutine.Instance().Go(slf.logMount)
	coroutine.Instance().Go(slf.logWithDefault)

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

func (slf *DefaultStart) logMount(args []interface{}) {
	slf.sysLogger.Mount()
}

func (slf *DefaultStart) logWithDefault(args []interface{}) {
	logger.WithDefault(slf.sysLogger)
}

/*TODO: 旧代码备份
coLimit := util.GetArgInt("colimit", util.MCCOPOOLDEFLIMIT)
coMax := util.GetArgInt("comax", util.MCCOPOOLDEFMAX)
coMin := util.GetArgInt("comin", util.MCCOPOOLDEFMIN)

logPath := util.GetArgString("logPath", "")
logLevl := util.GetArgInt("logLevel", int(logger.TRACELEVEL))
logSize := util.GetArgInt("logSize", 1024)
virDir := util.GetArgString("dir", "")

//初始化协程池
util.InitCoPool(coLimit, coMax, coMin)
//设置系统日志
s.sysLogger = logger.New(func() logger.Logger {
	l := logger.LogContext{FilName: logPath,
		LogHandle:  logrus.New(),
		LogMailbox: make(chan logger.Event, logSize),
		LogStop:    make(chan struct{})}

	l.LogHandle.SetLevel(logrus.Level(logLevl))

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	if runtime.GOOS == "windows" {
		formatter.DisableColors = true
	} else {
		formatter.SetColorScheme(&prefixed.ColorScheme{
			PrefixStyle: "blue+b"})
	}
	l.LogHandle.SetFormatter(formatter)
	l.Redirect()
	return &l
})
//---------------------
go s.sysLogger.Mount()
//---------------------
logger.WithDefault(s.sysLogger)
// 设置虚拟文件系统根目录
if virDir == "" {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	files.WithRootPath(dir)
} else {
	files.WithRootPath(virDir)
}*/
