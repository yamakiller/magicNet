package debug

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// Trace : 跟踪异常及崩溃记录到core文件中
func Trace() {
	if err := recover(); err != nil {
		saveCore(err)
		os.Exit(0)
	}
}

func saveCore(err interface{}) {
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-02015-04-05")
	fileName := "core-" + formatTimeStr + ". cre"

	f, err := os.Create(fileName)
	defer f.Close()

	f.WriteString(fmt.Sprintln(err))
	f.WriteString(stack())
	f.Sync()
}

func stack() string {
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
