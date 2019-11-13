package test

import (
	"testing"
	"time"

	"github.com/yamakiller/magicNet/engine/monitor"
	"github.com/yamakiller/magicNet/timer"
)

func TestTimer(t *testing.T) {
	timer.StartService()
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
	}
}
