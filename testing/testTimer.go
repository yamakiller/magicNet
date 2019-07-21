package testing

import (
	"time"

	"github.com/yamakiller/magicNet/engine/monitor"
	"github.com/yamakiller/magicNet/timer"
)

func TestTimer() {
	timer.StartService()
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
	}
}
