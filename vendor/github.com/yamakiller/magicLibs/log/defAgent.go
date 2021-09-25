package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

//DefaultAgent 默认日志代理
type DefaultAgent struct {
	_handle *logrus.Logger
}

//Close 关闭系统
func (slf *DefaultAgent) Close() {
	slf._handle = nil
}

//WithHandle 设置日志Handle
func (slf *DefaultAgent) WithHandle(handle interface{}) {
	slf._handle, _ = handle.(*logrus.Logger)
}

//Out 输出无等级日志
func (slf *DefaultAgent) Out(fmrt string, args ...interface{}) {
	slf._handle.Printf(fmrt, args...)
}

//Info 输出Info级日志
func (slf *DefaultAgent) Info(prefix string, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).Infoln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Infoln(fmt.Sprintf(fmrt, args...))
}

//Error 输出Error级日志
func (slf *DefaultAgent) Error(prefix string, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Errorln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Errorln(fmt.Sprintf(fmrt, args...))
}

//Debug 输出Debug级日志
func (slf *DefaultAgent) Debug(prefix string, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Debugln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Debugln(fmt.Sprintf(fmrt, args...))
}

//Warning 输出Warning级日志
func (slf *DefaultAgent) Warning(prefix, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Warningln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Warningln(fmt.Sprintf(fmrt, args...))
}

//Trace 输出Trace级日志
func (slf *DefaultAgent) Trace(prefix, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Traceln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Traceln(fmt.Sprintf(fmrt, args...))
}

//Fatal 输出Fatal级日志
func (slf *DefaultAgent) Fatal(prefix, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Fatalln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Fatalln(fmt.Sprintf(fmrt, args...))
}

//Panic 输出Panic级日志
func (slf *DefaultAgent) Panic(prefix, fmrt string, args ...interface{}) {
	if prefix != "" {
		slf._handle.WithFields(logrus.Fields{"prefix": prefix}).
			Panicln(fmt.Sprintf(fmrt, args...))
		return
	}

	slf._handle.Panicln(fmt.Sprintf(fmrt, args...))
}
