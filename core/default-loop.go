package core

import (
	"time"

	"github.com/yamakiller/magicNet/engine/monitor"
)

// DefaultLoop : 默认主循环体
type DefaultLoop struct {
}

// Wait : 主要循环体检测函数
func (dp *DefaultLoop) Wait() int {
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
		return 0
	}
	return -1
}
