package util

import (
	"sync/atomic"
)

// SpinLock : spin lock
type SpinLock struct {
	kernel uint32
}

// Trylock : try lock if unlock return false
func (slf *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&slf.kernel, 0, 1)
}

// Lock : locking
func (slf *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&slf.kernel, 0, 1) {
	}
}

// Unlock : unlocking
func (slf *SpinLock) Unlock() {
	atomic.StoreUint32(&slf.kernel, 0)
}
