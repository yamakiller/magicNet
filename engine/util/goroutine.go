package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

//GetCurrentGoroutineID : 获取当前协程的 ID
func GetCurrentGoroutineID() int {
	defer func() {
		if err := recover(); err != nil {
			//TODO: 需要发出警告
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	Assert(err == nil, fmt.Sprintf("cannot get goroutine id: %v", err))
	return id
}
