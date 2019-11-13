package util

import (
	"fmt"
	"runtime"
	"strings"
)

//Assert desc
//@method Assert desc: Assert boolean and output error message
//@param (bool) false assert
//@param (string) error message
func Assert(isAs bool, errMsg string) {
	if !isAs {
		panic(errMsg)
	}
}

//AssertEmpty desc
//@method AssertEmtpy desc: Assert Nil and output an error message
//@param (interface{}) is null assert
//@param (string) error message
func AssertEmpty(isNull interface{}, errMsg string) {
	if isNull == nil {
		panic(errMsg)
	}
}

//GetStack desc
//@method GetStack desc: Return current stack information
//@return (string)
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
