package testing

import (
	"magicNet/engine/monitor"
	"magicNet/timer"
	"time"
)

func TestTimer() {
	timer.StartService()
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
	}
}
