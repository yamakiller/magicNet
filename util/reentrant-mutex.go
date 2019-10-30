package util

import (
	"sync"
)

// ReMutex : Reentrant lock
type ReMutex struct {
	mutex *sync.Mutex
	owner int
	count int
}

// Association Sync lock association reentrant lock
func (slf *ReMutex) Association(m *sync.Mutex) {
	slf.mutex = m
}

// Lock : lock
func (slf *ReMutex) Lock() {
	me := GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return
	}

	slf.mutex.Lock()
}

// Unlock : unlock
func (slf *ReMutex) Unlock() {
	Assert(slf.owner == GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf.count > 0 {
		slf.count--
	} else {
		slf.mutex.Unlock()
	}
}
