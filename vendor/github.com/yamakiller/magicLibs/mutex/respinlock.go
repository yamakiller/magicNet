package mutex

import "github.com/yamakiller/magicLibs/util"

//ReSpinLock doc
//@Struct ReSpinLock @Summary Re-entrant spin lock
type ReSpinLock struct {
	_mutex *SpinLock
	_owner int
	_count int
}

//Width doc
//@Method Width doc : Spinlock association reentrant spin lock
//@Param (*SpinLock) width spinlock object
func (slf *ReSpinLock) Width(m *SpinLock) {
	slf._mutex = m
}

//Trylock doc
//@Method Trylock doc : Try to lock if you fail to get the lock return failure will not try again
//@Return (bool)
func (slf *ReSpinLock) Trylock() bool {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return true
	}

	return slf._mutex.Trylock()
}

//Lock doc
//@Method Lock @Summary locking
func (slf *ReSpinLock) Lock() {
	me := util.GetCurrentGoroutineID()
	if slf._owner == me {
		slf._count++
		return
	}

	slf._mutex.Lock()
}

//Unlock doc
//@Method Unlock @Summary unlocking
func (slf *ReSpinLock) Unlock() {
	util.Assert(slf._owner == util.GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf._count > 0 {
		slf._count--
	} else {
		slf._mutex.Unlock()
	}
}
