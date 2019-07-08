package actor

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ErrTimeout  : 定义默认超时错误
var ErrTimeout = errors.New("future: timeout")

// NewFuture :  创建一个Future
func NewFuture(d time.Duration) *Future {
	ref := &futureProcess{Future{cond: sync.NewCond(&sync.Mutex{})}}
	pid := &PID{p: (*Process)(unsafe.Pointer(ref))}
	GlobalRegistry.Register(pid)

	ref.pid = pid
	if d >= 0 {
		tp := time.AfterFunc(d, func() {
			ref.err = ErrTimeout
			ref.Stop(pid)
		})
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ref.t)), unsafe.Pointer(tp))
	}
	return &ref.Future
}

// Future : Future 对象
type Future struct {
	pid  *PID
	cond *sync.Cond

	done        bool
	result      interface{}
	err         error
	t           *time.Timer //是否需要修改
	pipes       []*PID
	completions []func(res interface{}, err error)
}

// PID : 获取对象的PID
func (f *Future) PID() *PID {
	return f.pid
}

//PipeTo : 多个PID 对象关联到Future
func (f *Future) PipeTo(pids ...*PID) {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	f.pipes = append(f.pipes, pids...)

	if f.done {
		f.sendToPipes()
	}
}

func (f *Future) sendToPipes() {
	if f.pipes == nil {
		return
	}

	var m interface{}
	if f.err != nil {
		m = f.err
	} else {
		m = f.result
	}

	for _, pid := range f.pipes {
		pid.sendUsrMessage(m)
	}
	f.pipes = nil
}

func (f *Future) wait() {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	for !f.done {
		f.cond.Wait()
	}
}

// Result : 获取结果
func (f *Future) Result() (interface{}, error) {
	f.wait()
	return f.result, f.err
}

// Wait : 等待结果
func (f *Future) Wait() error {
	f.wait()
	return f.err
}

func (f *Future) continueWith(continuation func(res interface{}, err error)) {
	f.cond.L.Lock()
	defer f.cond.L.Unlock()
	if f.done {
		continuation(f.result, f.err)
	} else {
		f.completions = append(f.completions, continuation)
	}
}

type futureProcess struct {
	Future
}

func (fp *futureProcess) SendUsrMessage(pid *PID, message interface{}) {
	_, msg, _ := UnWrapPack(message)
	fp.result = msg
	fp.Stop(pid)
}

func (fp *futureProcess) SendSysMessage(pid *PID, message interface{}) {
	fp.result = message
	fp.Stop(pid)
}

func (fp *futureProcess) Stop(pid *PID) {
	fp.cond.L.Lock()
	if fp.done {
		fp.cond.L.Unlock()
		return
	}

	fp.done = true
	tp := (*time.Timer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&fp.t))))
	if tp != nil {
		tp.Stop()
	}
	GlobalRegistry.UnRegister(pid)

	fp.sendToPipes()
	fp.runCompletions()
	fp.cond.L.Unlock()
	fp.cond.Signal()
}

func (f *Future) runCompletions() {
	if f.completions == nil {
		return
	}

	for _, c := range f.completions {
		c(f.result, f.err)
	}
	f.completions = nil
}
