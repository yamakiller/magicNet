package util

// ReSpinLock : 可重入自旋锁
type ReSpinLock struct {
	mutex *SpinLock
	owner int
	count int
}

// Association : 自旋锁关联可重入自旋锁
func (re *ReSpinLock) Association(m *SpinLock) {
	re.mutex = m
}

// TryLock : 尝试加锁如果未获得锁返回失败不会反复尝试
func (re *ReSpinLock) TryLock() bool {
	me := GetCurrentGoroutineID()
	if re.owner == me {
		re.count++
		return true
	}

	return re.mutex.Trylock()
}

// Lock : 加锁
func (re *ReSpinLock) Lock() {
	me := GetCurrentGoroutineID()
	if re.owner == me {
		re.count++
		return
	}

	re.mutex.Lock()
}

// Unlock : 解锁
func (re *ReSpinLock) Unlock() {
	Assert(re.owner == GetCurrentGoroutineID(), "illegalMonitorStateError")
	if re.count > 0 {
		re.count--
	} else {
		re.mutex.Unlock()
	}
}
