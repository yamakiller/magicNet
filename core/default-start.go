package core

import (
	"os"
	"runtime"

	"github.com/yamakiller/magicNet/engine/util"

	"github.com/yamakiller/magicNet/engine/files"
	"github.com/yamakiller/magicNet/engine/logger"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// DefaultStart :默认启动器
type DefaultStart struct {
	sysLogger logger.Logger
}

// Init : 初始化系统
func (s *DefaultStart) Init() error {

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
	}

	return nil
}

// Destory : 销毁处理
func (s *DefaultStart) Destory() {
	if s.sysLogger != nil {
		s.sysLogger.Close()
	}

	files.Close()
	util.DestoryCoPool()
}
