package core

import (
	"time"

	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/engine/monitor"
)

//DefaultWait desc
//@struct DefaultWait desc: Default main loop body
//@member (*util.SignalWatch) watch exit signal
type DefaultWait struct {
	_sw *util.SignalWatch
}

//Enter desc
//@method Enter desc: enter system
func (slf *DefaultWait) Enter() {

	if !monitor.IsShutdown() {
		return
	}

	slf._sw = &util.SignalWatch{}
	slf._sw.Initial(slf.shutdown)
	slf._sw.Watch()
}

//Wait desc
//@method Wait desc: Waiting for the system to be terminated
//@return (int) 0:continue -1:eixt
func (slf *DefaultWait) Wait() int {
	if !monitor.IsShutdown() {
		time.Sleep(time.Millisecond * 1000)
		return 0
	}

	slf._sw.Wait()

	return -1
}

func (slf *DefaultWait) shutdown() {
	monitor.Shutdown()
}
