package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Logger : 日志模块接口
type Logger interface {
	run() int
	exit()
	redirect()
	write(msg *event)
	getPrefix(owner uint32) string

	Mount()
	Close()
	Error(owner uint32, fmrt string, args ...interface{})
	Info(owner uint32, fmrt string, args ...interface{})
	Warning(owner uint32, fmrt string, args ...interface{})
	Panic(owner uint32, fmrt string, args ...interface{})
	Fatal(owner uint32, fmrt string, args ...interface{})
	Debug(owner uint32, fmrt string, args ...interface{})
	Trace(owner uint32, fmrt string, args ...interface{})
}

// LogContext : 日志对象
type LogContext struct {
	filName    string
	filHandle  *os.File
	logLevel   logrus.Level
	logHandle  *logrus.Logger
	logMailNum int32
	logMailbox chan event
	logStop    chan struct{}
	logWait    sync.WaitGroup
}

// MakeLogger : 日志制作器
type MakeLogger func() Logger

var (
	defaultLevel      = logrus.PanicLevel
	defaultSize       = 512
	defaultFile       = ""
	defaultMakeLogger = func() Logger {
		l := LogContext{filName: defaultFile,
			logHandle:  logrus.New(),
			logMailbox: make(chan event, defaultSize),
			logStop:    make(chan struct{})}

		l.logHandle.SetOutput(os.Stdout)
		l.logHandle.SetLevel(l.logLevel)
		formatter := new(prefixed.TextFormatter)
		formatter.FullTimestamp = true
		if runtime.GOOS == "windows" {
			formatter.DisableColors = true
		} else {
			formatter.SetColorScheme(&prefixed.ColorScheme{
				PrefixStyle: "blue+b"})
		}
		l.logHandle.SetFormatter(formatter)
		l.redirect()
		return &l
	}

	defaultHandle Logger
)

// New : 创建日志对象
func New(maker MakeLogger) *Logger {

	if maker == nil {
		r := defaultMakeLogger()
		return &r
	}

	r := maker()
	return &r
}

// WithDefault : 关联
func WithDefault(log Logger) {
	defaultHandle = log
}

// Error : 输出错误日志
func Error(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Error(owner, fmrt, args...)
}

// Info : 输出信息日志
func Info(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Info(owner, fmrt, args...)
}

// Warning : 输出警告日志
func Warning(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Warning(owner, fmrt, args...)
}

// Panic : 输出程序崩溃日志
func Panic(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Panic(owner, fmrt, args...)
}

// Fatal : 输出严重错误日志
func Fatal(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Fatal(owner, fmrt, args...)
}

// Debug : 输出Debug日志
func Debug(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Debug(owner, fmrt, args...)
}

// Trace : 输出跟踪日志
func Trace(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Trace(owner, fmrt, args...)
}

func (log *LogContext) run() int {
	select {
	case <-log.logStop:
		return -1
	case msg := <-log.logMailbox:
		log.write(&msg)
		atomic.AddInt32(&log.logMailNum, -1)
		return 0
	}
}

func (log *LogContext) exit() {
	log.logWait.Done()
}

func (log *LogContext) write(msg *event) {
	switch msg.level {
	case uint32(logrus.ErrorLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Errorln(msg.message)
	case uint32(logrus.InfoLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Infoln(msg.message)
	case uint32(logrus.TraceLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Traceln(msg.message)
	case uint32(logrus.DebugLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Debugln(msg.message)
	case uint32(logrus.WarnLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Warningln(msg.message)
	case uint32(logrus.FatalLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Fatalln(msg.message)
	case uint32(logrus.PanicLevel):
		log.logHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Panicln(msg.message)
	}
}

func (log *LogContext) redirect() {
	if strings.Compare(log.filName, "") == 0 {
		return
	}

	f, err := os.OpenFile(log.filName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	log.filHandle = f
	log.logHandle.SetOutput(f)
}

func (log *LogContext) getPrefix(owner uint32) string {
	if owner == 0 {
		return "[main]"
	}
	return fmt.Sprintf("%08x", owner)
}

func (log *LogContext) push(data event) {
	log.logMailbox <- data
	atomic.AddInt32(&log.logMailNum, 1)
}

// Mount : 挂载日志模块
func (log *LogContext) Mount() {
	log.logWait.Add(1)
	go func(log Logger) {
		for {
			if log.run() != 0 {
				break
			}
		}
		log.exit()
	}(log)
}

// Close : 关闭日志系统
func (log *LogContext) Close() {
	for {
		if atomic.LoadInt32(&log.logMailNum) > 0 {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}

	close(log.logStop)
	log.logWait.Wait()
	close(log.logMailbox)
}

// Error : 输出错误日志
func (log *LogContext) Error(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.ErrorLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})

}

// Info : 输出信息日志
func (log *LogContext) Info(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.InfoLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

// Warning : 输出警告日志
func (log *LogContext) Warning(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.WarnLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

// Panic : 输出程序崩溃日志
func (log *LogContext) Panic(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.PanicLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

// Fatal : 输出严重错误日志
func (log *LogContext) Fatal(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.FatalLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

// Debug : 输出Debug日志
func (log *LogContext) Debug(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.DebugLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

// Trace : 输出跟踪日志
func (log *LogContext) Trace(owner uint32, fmrt string, args ...interface{}) {
	log.push(event{level: uint32(logrus.TraceLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}
