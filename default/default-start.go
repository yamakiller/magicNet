package frame

import (
	"flag"
	"magicNet/engine/files"
	"magicNet/engine/logger"
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// DefaultStart :默认启动器
type DefaultStart struct {
	sysLogger logger.Logger
}

// Init : 初始化系统
func (s *DefaultStart) Init() error {
	logPath := flag.String("logPath", "", "log file path")
	logLevl := flag.Int("logLevel", int(logger.PANICLEVEL), "log level")
	logSize := flag.Int("logSize", 1024, "log mailbox size")
	virDir := flag.String("v", "", "virtual root directory")

	//设置系统日志
	s.sysLogger = logger.New(func() logger.Logger {
		l := logger.LogContext{FilName: *logPath,
			LogHandle:  logrus.New(),
			LogMailbox: make(chan logger.Event, *logSize),
			LogStop:    make(chan struct{})}
		l.LogHandle.SetLevel(logrus.Level(*logLevl))
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
	logger.WithDefault(s.sysLogger)
	// 设置虚拟文件系统根目录
	if *virDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		files.WithRootPath(dir)
	} else {
		files.WithRootPath(*virDir)
	}

	return nil
}

// Destory : 销毁处理
func (s *DefaultStart) Destory() {
	if s.sysLogger != nil {
		s.sysLogger.Close()
	}

	files.Close()
}
