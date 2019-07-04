package util

import (
  "sync/atomic"
)

type SpinLock struct {
  kernel uint32
}

func (sl *SpinLock) Trylock() bool {
  return atomic.CompareAndSwapUint32(&sl.kernel, 0, 1)
}

func (sl *SpinLock) Lock() {
    for !atomic.CompareAndSwapUint32(&sl.kernel, 0, 1) {
    }
}

func (sl *SpinLock) Unlock() {
  atomic.StoreUint32(&sl.kernel, 0)
}
