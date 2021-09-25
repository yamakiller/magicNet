package mutex

import (
	"sync/atomic"
	"time"
)

//SpinLock doc
//@Struct SpinLick @Summary spin lock
//@Member (uint32)
type SpinLock struct {
	Deplay  time.Duration
	Check   int
	_kernel uint32
}

//Trylock doc
//@Method Trylock @Summary try lock if unlock return false
//@Return (bool)
func (slf *SpinLock) Trylock() bool {
	return atomic.CompareAndSwapUint32(&slf._kernel, 0, 1)
}

//Lock doc
//@Method Lock @Summary locking
func (slf *SpinLock) Lock() {
	tmpCheck := 0
	tmpDeplay := slf.Deplay
	for !atomic.CompareAndSwapUint32(&slf._kernel, 0, 1) {
		if tmpDeplay > 0 {
			tmpCheck++
			if tmpCheck < slf.Check {
				continue
			}
			tmpDeplay *= 5
			if max := time.Duration(500) * time.Millisecond; tmpDeplay > max {
				tmpDeplay = max
			}
			tmpCheck = 0
			time.Sleep(tmpDeplay)
		}
	}
}

//Unlock doc
//@Method Unlock @Summary unlocking
func (slf *SpinLock) Unlock() {
	atomic.StoreUint32(&slf._kernel, 0)
}
