package debug

import (
	"fmt"
	"magicNet/core/version"
	"os"
	"runtime"
	"runtime/trace"
	"strings"
	"time"
)

// TraceDebug : Debug 跟踪器
type TraceDebug struct {
}

// Start : 启动DEBUG跟踪器
func (t *TraceDebug) Start() {
	if strings.TrimSpace(strings.ToLower(version.Build)) == "debug" {
		trace.Start(os.Stderr)
	}
}

// Stop : 停止DEBUG跟踪器
func (t *TraceDebug) Stop() {
	if strings.TrimSpace(strings.ToLower(version.Build)) == "release" {
		if err := recover(); err != nil {
			t.saveCore(err)
			os.Exit(0)
		}
	}

	if strings.TrimSpace(strings.ToLower(version.Build)) == "debug" {
		trace.Stop()
	}
}

func (t *TraceDebug) saveCore(err interface{}) {
	timeUnix := time.Now().Unix()
	formatTimeStr := time.Unix(timeUnix, 0).Format("2006-01-04-05")
	fileName := "core-" + formatTimeStr + ". cre"

	f, err := os.Create(fileName)
	defer f.Close()

	f.WriteString(fmt.Sprintln(err))
	f.WriteString(t.stack())
	f.Sync()
}

func (t *TraceDebug) stack() string {
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
