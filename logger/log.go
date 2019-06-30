package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logf *os.File

// Init is Initialization log module
func Init(logLevel string) {
	logf = nil
	logrus.SetOutput(os.Stdout)
	configLevel(logLevel)
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle: "blue+b",
	})

	logrus.SetFormatter(formatter)
}

// Redirect : 重定向
func Redirect(filename string) {
	if strings.Compare(filename, "") != 0 {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			Error(0, "redirect log file fail:%s", filename)
		} else {
			logf = f
			logrus.SetOutput(logf)
		}
	}
}

// Destory : 销毁日志模块
func Destory() {
	if logf == nil {
		return
	}

	logf.Close()
}

func configLevel(lvl string) {
	lv, err := logrus.ParseLevel(lvl)
	if err != nil {
		return
	}

	logrus.SetLevel(lv)
}

func wrap(fmrt string) string {
	return fmrt
}

func prefix(owner uint32) string {
	if owner == 0 {
		return "[main]"
	} else {
		return fmt.Sprintf("%08x", owner)
	}
}

// Error write error message
func Error(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Errorf(wrap(fmrt), args...)
}

// Info write message
func Info(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Infof(wrap(fmrt), args...)
}

// Warning write warning message
func Warning(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Warningf(wrap(fmrt), args...)
}

//Panic write serious error message
func Panic(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Panicf(wrap(fmrt), args...)
}

//Fatal write fatal error message
func Fatal(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Fatalf(wrap(fmrt), args...)
}

//Trace write trace message
func Trace(owner uint32, fmrt string, args ...interface{}) {
	logrus.WithFields(logrus.Fields{"prefix": prefix(owner)}).Tracef(wrap(fmrt), args...)
}
