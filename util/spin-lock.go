package util

import (
	"sync/atomic"
)

// SpinLock : 自旋锁
type SpinLock struct {
	kernel uint32
}

// Trylock : 尝试竞争锁如果未竞争到返回false
func (sl *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&sl.kernel, 0, 1)
}

// Lock : 加锁
func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapUint32(&sl.kernel, 0, 1) {
	}
}

// Unlock : 解锁
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.kernel, 0)
}
