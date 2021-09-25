package actors

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/yamakiller/magicLibs/actors/messages"
	"github.com/yamakiller/magicLibs/util"
)

const (
	idle int32 = iota
	running
)

type mailbox struct {
	_usrMailbox      *queue
	_sysMailbox      *queue
	_usrMessages     int32
	_sysMessages     int32
	_schedulerStatus int32
	_dispatcher      *dispatcher
	_invoker         invoker
	_suspended       int32
}

func (slf *mailbox) postSysMessage(msg interface{}) {
	slf._sysMailbox.Push(msg)
	atomic.AddInt32(&slf._sysMessages, 1)
	slf.schedule()
}

func (slf *mailbox) postUsrMessage(msg interface{}) {
	slf._usrMailbox.Push(msg)
	atomic.AddInt32(&slf._usrMessages, 1)
	slf.schedule()
}

func (slf *mailbox) schedule() {
	if atomic.CompareAndSwapInt32(&slf._schedulerStatus, idle, running) {
		slf._dispatcher.Schedule(slf.processMessages)
	}
}

func (slf *mailbox) run() {
	var msg interface{}
	//Begin=>致命性异常处理
	defer func() {
		if r := recover(); r != nil {
			slf._invoker.escalateFailure(r, msg)
		}
	}()
	//End=>致命性异常处理结束
	i := 0
	for {
		if i > 0 {
			i = 0
			runtime.Gosched()
		}

		i++
		if msg, _ = slf._sysMailbox.Pop(); msg != nil {
			atomic.AddInt32(&slf._sysMessages, -1)
			switch msg.(type) {
			case messages.Suspend:
				atomic.StoreInt32(&slf._suspended, 1)
			case messages.Resume:
				atomic.StoreInt32(&slf._suspended, 0)
			default:
				slf._invoker.invokeSysMessage(msg)
			}
			continue
		}

		if atomic.LoadInt32(&slf._suspended) == 1 {
			return
		}

		if msg, _ = slf._usrMailbox.Pop(); msg != nil {
			atomic.AddInt32(&slf._usrMessages, -1)
			slf._invoker.invokeUsrMessage(msg)
		} else {
			return
		}

	}
}

func (slf *mailbox) processMessages([]interface{}) {
process_lable:
	slf.run()

	atomic.StoreInt32(&slf._schedulerStatus, idle)
	sys := atomic.LoadInt32(&slf._sysMessages)
	usr := atomic.LoadInt32(&slf._usrMessages)

	if sys > 0 || (atomic.LoadInt32(&slf._suspended) == 0 && usr > 0) {
		if atomic.CompareAndSwapInt32(&slf._schedulerStatus, idle, running) {
			goto process_lable
		}
	}
}

func spawnQueue(cap int) *queue {
	return &queue{
		_cap:               cap,
		_overloadThreshold: cap * 2,
		_buffer:            make([]interface{}, cap),
	}
}

type queue struct {
	_cap  int
	_head int
	_tail int

	_overload          int
	_overloadThreshold int

	_buffer []interface{}
	_sync   sync.Mutex
}

//Push Insert an object
//@Param (interface{}) item
func (slf *queue) Push(item interface{}) {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	slf.unpush(item)
}

//Pop doc
//@Method Pop @Summary Take an object, If empty return nil
//@Return (interface{}) return object
//@Return (bool)
func (slf *queue) Pop() (interface{}, bool) {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	return slf.unpop()
}

//Overload Detecting queues exceeding the limit [mainly used for warning records]
//@Return (int)
func (slf *queue) Overload() int {
	if slf._overload != 0 {
		overload := slf._overload
		slf._overload = 0
		return overload
	}
	return 0
}

//Length Length of the queue
//@Return (int) length
func (slf *queue) Length() int {
	var (
		head int
		tail int
		cap  int
	)
	slf._sync.Lock()
	head = slf._head
	tail = slf._tail
	cap = slf._cap
	slf._sync.Unlock()

	if head <= tail {
		return tail - head
	}
	return tail + cap - head
}

func (slf *queue) unpush(item interface{}) {
	util.AssertEmpty(item, "error push is nil")
	slf._buffer[slf._tail] = item
	slf._tail++
	if slf._tail >= slf._cap {
		slf._tail = 0
	}

	if slf._head == slf._tail {
		slf.expand()
	}
}

func (slf *queue) unpop() (interface{}, bool) {
	var resultSucces bool
	var result interface{}
	if slf._head != slf._tail {
		resultSucces = true
		result = slf._buffer[slf._head]
		slf._buffer[slf._head] = nil
		slf._head++
		if slf._head >= slf._cap {
			slf._head = 0
		}

		length := slf._tail - slf._tail
		if length < 0 {
			length += slf._cap
		}
		for length > slf._overloadThreshold {
			slf._overload = length
			slf._overloadThreshold *= 2
		}
	}

	return result, resultSucces
}

func (slf *queue) expand() {
	newBuff := make([]interface{}, slf._cap*2)
	for i := 0; i < slf._cap; i++ {
		newBuff[i] = slf._buffer[(slf._head+i)%slf._cap]
	}

	slf._head = 0
	slf._tail = slf._cap
	slf._cap *= 2

	slf._buffer = newBuff
}
