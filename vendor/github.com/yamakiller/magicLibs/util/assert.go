package util

import (
	"fmt"
	"runtime"
	"strings"
)

//Assert Assert boolean and output error message
//@Method Assert
//@Param (bool) false assert
//@Param (string) error message
func Assert(isAs bool, errMsg string) {
	if !isAs {
		_, file, inline, ok := runtime.Caller(2)
		panic(fmt.Sprintf("%s %d %v\n%s", file, inline, ok, errMsg))
	}
}

//AssertError Assert error is null
func AssertError(err error) {
	if err != nil {
		panic(err)
	}
}

//AssertEmpty Assert Nil and output an error message
//@Method AssertEmtpy
//@Param (interface{}) is null assert
//@Param (string) error message
func AssertEmpty(isNull interface{}, errMsg string) {
	if isNull == nil {
		_, file, inline, ok := runtime.Caller(2)
		panic(fmt.Sprintf("%s %d %v\n%s", file, inline, ok, errMsg))
	}
}

//GetStack Return current stack information
//@Method GetStack
//@Return (string)
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
