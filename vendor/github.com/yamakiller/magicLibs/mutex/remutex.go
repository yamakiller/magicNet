package mutex

import (
	"sync"

	"github.com/yamakiller/magicLibs/util"
)

//ReMutex doc
//@Struct ReMutex @Summary Reentrant mutex
type ReMutex struct {
	_mutex *sync.Mutex
	_owner int
	_count int
}

//Width doc
//@Method Width @Summary Sync lock association reentrant lock
//@Param (*sync.Mutex) mutex object
func (slf *ReMutex) Width(m *sync.Mutex) {
	slf._mutex = m
}

//Lock doc
//@Method Lock @Summary locking
func (slf *ReMutex) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return
	}

	slf._mutex.Lock()
}

//Unlock doc
//@Method Unlock doc : unlocking
func (slf *ReMutex) Unlock() {
	util.Assert(slf._owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf._count > 0 {
		slf._count--
	} else {
		slf._mutex.Unlock()
	}
}
