package util

// ReSpinLock : Re-entrant spin lock
type ReSpinLock struct {
	mutex *SpinLock
	owner int
	count int
}

// Association : Spinlock association reentrant spin lock
func (slf *ReSpinLock) Association(m *SpinLock) {
	slf.mutex = m
}

// TryLock : Try to lock if you fail to get the lock return failure will not try again
func (slf *ReSpinLock) TryLock() bool {
	me := GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return true
	}

	return slf.mutex.Trylock()
}

// Lock : lock
func (slf *ReSpinLock) Lock() {
	me := GetCurrentGoroutineID()
	if slf.owner == me {
		slf.count++
		return
	}

	slf.mutex.Lock()
}

// Unlock : unlock
func (slf *ReSpinLock) Unlock() {
	Assert(slf.owner == GetCurrentGoroutineID(), "illegalMonitorStateError")
	if slf.count > 0 {
		slf.count--
	} else {
		slf.mutex.Unlock()
	}
}
