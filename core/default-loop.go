package core

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yamakiller/magicNet/engine/monitor"
)

// DefaultLoop : 默认主循环体
type DefaultLoop struct {
	c chan os.Signal
	e sync.WaitGroup
}

//EnterLoop desc
//@method EnterLoop desc: enter loop Pretreatment
func (slf *DefaultLoop) EnterLoop() {
	slf.c = make(chan os.Signal)
	slf.e = sync.WaitGroup{}

	if !monitor.IsShutdown() {
		return
	}

	slf.e.Add(1)
	signal.Notify(slf.c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		defer slf.e.Done()
		for s := range slf.c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				monitor.Shutdown()
				break
			default:
				break
			}
		}

	}()
}

// Wait desc
//@method Wait desc: Waiting for the system to be terminated
//@return (int) 0:continue -1:eixt
func (slf *DefaultLoop) Wait() int {
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
		return 0
	}

	slf.e.Wait()

	return -1
}
