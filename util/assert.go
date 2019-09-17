package util

import (
	"fmt"
	"runtime"
	"strings"
)

// Assert :  断言Bool并输出错误信息
func Assert(isAs bool, errMsg string) {
	if !isAs {
		panic(errMsg)
	}
}

// AssertEmpty : 断言Nil并输出错误信息
func AssertEmpty(isNull interface{}, errMsg string) {
	if isNull == nil {
		panic(errMsg)
	}
}

// GetStack : 获取当前堆栈信息
func GetStack() string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(4, pc[:])
	callers := pc[:n]
	frames := runtime.CallersFrames(callers)
	for {
		frame, more := frames.Next()
		file = frame.File
		line = frame.Line
		name = frame.Function
		if !strings.HasPrefix(name, "runtime.") || !more {
			break
		}
	}

	var str string
	switch {
	case name != "":
		str = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		str = fmt.Sprintf("%v:%v", file, line)
	default:
		str = fmt.Sprintf("pc:%x", pc)
	}
	return "stacktrace:\n" + str
}
