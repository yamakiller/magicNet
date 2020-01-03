package boot

import (
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/core/debug"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/core/version"
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

	//start-up logger
	logPath := args.Instance().GetString("-log", "")
	logSize := args.Instance().GetInt("-logSize", 128)
	logLevel := logger.DEBUGLEVEL

	release := args.Instance().GetBoolean("-release", false)
	if release {
		logLevel = logger.INFOLEVEL
	}

	log := logger.New(func() logger.Logger {
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

	logger.WithDefault(log)
	log.Mount()

	var space string
	for i := 0; i < ((54 - len(version.BuildName)) / 2); i++ {
		space += " "
	}

	logger.Info(0, "                     _          __     _         ____")
	logger.Info(0, "   /\\/\\   __ _  __ _(_) ___  /\\ \\ \\___| |_      /\\___\\")
	logger.Info(0, "  /    \\ / _` |/ _` | |/ __|/  \\/ / _ \\ __|    /\\ \\___\\")
	logger.Info(0, " / /\\/\\ \\ (_| | (_| | | (__/ /\\  /  __/ |_     \\ \\/ / /")
	logger.Info(0, " \\/    \\/\\__,_|\\__, |_|\\___\\_\\ \\/ \\___|\\__|     \\/_/_/")
	logger.Info(0, " ::magic net:: |___/(v%s %s %s)", version.BuildVersion, "DEBUG", version.BuildTime)
	logger.Info(0, " ::%s %s", version.CommitID, util.TimeNowFormat())
	logger.Info(0, "--------------------------------------------------------")
	logger.Info(0, "| %s%s%s|", space, version.BuildName, space)
	logger.Info(0, "--------------------------------------------------------")

	fme := f()
	if err := fme.Initial(); err != nil {
		logger.Error(0, "%+v", err)
		goto end
	}

	if err := fme.InitService(); err != nil {
		logger.Error(0, "%+v", err)
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
	log.Close()
}
