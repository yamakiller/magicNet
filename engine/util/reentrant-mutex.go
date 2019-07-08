package util

import (
	"sync"
)

// ReMutex : 可重入锁
type ReMutex struct {
	mutex *sync.Mutex
	owner int
	count int
}

// Association 同步锁关联可重入锁
func (re *ReMutex) Association(m *sync.Mutex) {
	re.mutex = m
}

// Lock : 加锁
func (re *ReMutex) Lock() {
	me := GetCurrentGoroutineID()
	if re.owner == me {
		re.count++
		return
	}

	re.mutex.Lock()
}

// Unlock : 解锁
func (re *ReMutex) Unlock() {
	Assert(re.owner == GetCurrentGoroutineID(), "illegalMonitorStateError")
	if re.count > 0 {
		re.count--
	} else {
		re.mutex.Unlock()
	}
}
