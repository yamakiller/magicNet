package logger

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

const (
	// PANICLEVEL Crash log level
	PANICLEVEL uint32 = iota
	// FATALLEVEL Critical error log level
	FATALLEVEL
	// ERRORLEVEL Error log level
	ERRORLEVEL
	// WARNLEVEL  Warning log level
	WARNLEVEL
	// INFOLEVEL  General information log level
	INFOLEVEL
	// DEBUGLEVEL Debug log level
	DEBUGLEVEL
	// TRACELEVEL Trace log level
	TRACELEVEL
)

//Logger desc
//@interface Logger desc: Log module interface
type Logger interface {
	run() int
	exit()

	write(msg *Event)
	getPrefix(owner uint32) string

	Mount()
	Redirect()
	Close()
	Error(owner uint32, fmrt string, args ...interface{})
	Info(owner uint32, fmrt string, args ...interface{})
	Warning(owner uint32, fmrt string, args ...interface{})
	Panic(owner uint32, fmrt string, args ...interface{})
	Fatal(owner uint32, fmrt string, args ...interface{})
	Debug(owner uint32, fmrt string, args ...interface{})
	Trace(owner uint32, fmrt string, args ...interface{})
}

//LogContext desc
//@struct LogContext desc: Log context
type LogContext struct {
	FilName    string
	FilHandle  *os.File
	LogLevel   logrus.Level
	LogHandle  *logrus.Logger
	LogMailNum int32
	LogMailbox chan Event
	LogStop    chan struct{}
	LogWait    sync.WaitGroup
}

//MakeLogger desc
//@method MakeLogger desc: Log object maker
type MakeLogger func() Logger

var (
	defaultLevel      = logrus.PanicLevel
	defaultSize       = 512
	defaultFile       = ""
	defaultMakeLogger = func() Logger {
		l := LogContext{FilName: defaultFile,
			LogHandle:  logrus.New(),
			LogMailbox: make(chan Event, defaultSize),
			LogStop:    make(chan struct{})}
		l.LogHandle.SetLevel(l.LogLevel)
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
	}

	defaultHandle Logger
)

// New : 创建日志对象
func New(maker MakeLogger) Logger {

	if maker == nil {
		r := defaultMakeLogger()
		return r
	}

	r := maker()
	return r
}

//WithDefault desc
//@method WithDefault desc: Set the default log handle
//@param (Logger) logger object
func WithDefault(log Logger) {
	defaultHandle = log
}

//Error desc
//@method Error desc: Output error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Error(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Error(owner, fmrt, args...)
}

//Info desc
//@method Info desc: Output information log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Info(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Info(owner, fmrt, args...)
}

//Warning desc
//@method Warning desc: Output warning log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Warning(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Warning(owner, fmrt, args...)
}

//Panic desc
//@method Panic desc: Output program crash log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Panic(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Panic(owner, fmrt, args...)
}

//Fatal desc
//@method Fatal desc: Output critical error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Fatal(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Fatal(owner, fmrt, args...)
}

//Debug desc
//@method Debug desc: Output Debug log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Debug(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Debug(owner, fmrt, args...)
}

//Trace desc
//@method Trace desc: Output trace log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func Trace(owner uint32, fmrt string, args ...interface{}) {
	if defaultHandle == nil {
		return
	}
	defaultHandle.Trace(owner, fmrt, args...)
}

func (log *LogContext) run() int {
	select {
	case <-log.LogStop:
		return -1
	case msg := <-log.LogMailbox:
		log.write(&msg)
		atomic.AddInt32(&log.LogMailNum, -1)
		return 0
	}
}

func (log *LogContext) exit() {
	log.LogWait.Done()
}

func (log *LogContext) write(msg *Event) {
	switch msg.level {
	case uint32(logrus.ErrorLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Errorln(msg.message)
	case uint32(logrus.InfoLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Infoln(msg.message)
	case uint32(logrus.TraceLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Traceln(msg.message)
	case uint32(logrus.DebugLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Debugln(msg.message)
	case uint32(logrus.WarnLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Warningln(msg.message)
	case uint32(logrus.FatalLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Fatalln(msg.message)
	case uint32(logrus.PanicLevel):
		log.LogHandle.WithFields(logrus.Fields{"prefix": msg.prefix}).Panicln(msg.message)
	}
}

func (log *LogContext) getPrefix(owner uint32) string {
	if owner == 0 {
		return "[&main]"
	}
	return fmt.Sprintf("[&%08x]", owner)
}

func (log *LogContext) push(data Event) {
	select {
	case log.LogMailbox <- data:
	}

	atomic.AddInt32(&log.LogMailNum, 1)
}

//Redirect desc
//@method Redirect desc: Redirect log file
func (log *LogContext) Redirect() {
	if log.FilName == "" {
		log.LogHandle.SetOutput(os.Stdout)
		return
	}

	f, err := os.OpenFile(log.FilName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}
	log.FilHandle = f
	log.LogHandle.SetOutput(f)
}

//Mount desc
//@method Mount desc: Mount log module
func (log *LogContext) Mount() {
	//TODO:需要修改
	log.LogWait.Add(1)
	go func(log Logger) {
		for {
			if log.run() != 0 {
				break
			}
		}
		log.exit()
	}(log)
}

//Close desc
//@method Close desc: Turn off the logging system
func (log *LogContext) Close() {
	for {
		if atomic.LoadInt32(&log.LogMailNum) > 0 {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}

	close(log.LogStop)
	log.LogWait.Wait()
	close(log.LogMailbox)
	if log.FilHandle != nil {
		log.FilHandle.Close()
	}
}

//Error desc
//@method Error desc: Output error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Error(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.ErrorLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})

}

//Info desc
//@method Info desc: Output information log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Info(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.InfoLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

//Warning desc
//@method Warning desc: Output warning log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Warning(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.WarnLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

//Panic desc
//@method Panic desc: Output program crash log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Panic(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.PanicLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

//Fatal desc
//@method Fatal desc: Output critical error log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Fatal(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.FatalLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

//Debug desc
//@method Debug desc: Output Debug log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Debug(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.DebugLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}

//Trace desc
//@method Trace desc: Output trace log
//@param (int32) owner
//@param (string) format
//@param (...interface{}) args
func (log *LogContext) Trace(owner uint32, fmrt string, args ...interface{}) {
	log.push(Event{level: uint32(logrus.TraceLevel), prefix: log.getPrefix(owner), message: fmt.Sprintf(fmrt, args...)})
}
