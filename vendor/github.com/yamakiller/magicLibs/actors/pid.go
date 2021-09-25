package actors

import (
	"fmt"
	"sync/atomic"

	"github.com/yamakiller/magicLibs/mutex"
)

//PID 外部接口
type PID struct {
	ID      uint32
	_h      handle
	_parent *Core
	_sync   mutex.SpinLock
}

//ToString 返回ID字符串
func (slf *PID) ToString() string {
	return fmt.Sprintf(".%08x", slf.ID)
}

func (slf *PID) ref() handle {
	p := slf._h
	if p != nil {
		if l, ok := (p).(*actorHandle); ok && atomic.LoadInt32(&l._death) == 1 {
			slf._sync.Lock()
			slf._h = nil
			slf._sync.Unlock()
		} else {
			return p
		}
	}

	ref := slf._parent.getHandle(slf)
	if ref != nil {
		slf._sync.Lock()
		slf._h = ref
		slf._sync.Unlock()
	}
	return ref
}

//Post 发送消息
func (slf *PID) Post(message interface{}) {
	slf.postUsrMessage(message)
}

func (slf *PID) postUsrMessage(message interface{}) {
	ref := slf.ref()
	ref.postUsrMessage(slf, message)

	overload := ref.overloadUsrMessage()
	if overload > 0 {
		slf.Warning("user mailbox overload %d", overload)
	}
}

func (slf *PID) postSysMessage(message interface{}) {
	slf.ref().postSysMessage(slf, message)
}

//Info ...
func (slf *PID) Info(sfmt string, args ...interface{}) {
	slf._parent._log.Info(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Debug ...
func (slf *PID) Debug(sfmt string, args ...interface{}) {
	slf._parent._log.Debug(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Error ...
func (slf *PID) Error(sfmt string, args ...interface{}) {
	slf._parent._log.Error(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Warning ...
func (slf *PID) Warning(sfmt string, args ...interface{}) {
	slf._parent._log.Warning(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Fatal ...
func (slf *PID) Fatal(sfmt string, args ...interface{}) {
	slf._parent._log.Fatal(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Panic ...
func (slf *PID) Panic(sfmt string, args ...interface{}) {
	slf._parent._log.Panic(fmt.Sprintf("[%s]", slf.ToString()), sfmt, args...)
}

//Stop 停止
func (slf *PID) Stop() {
	slf.ref().Stop(slf)
}
