package frame

import (
	"flag"
	"magicNet/engine/files"
	"magicNet/engine/logger"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type defaultStart struct {
	sysLogger logger.Logger
}

func (s *defaultStart) Start() bool {
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
	logger.WithDefault(s.sysLogger)
	// 设置虚拟文件系统根目录
	if strings.Compare(*virDir, "") == 0 {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		files.WithRootPath(dir)
	} else {
		files.WithRootPath(*virDir)
	}

	return true
}

func (s *defaultStart) Shutdown() {
	if s.sysLogger != nil {
		s.sysLogger.Close()
	}

	files.Close()
}
