package util

type ReSpinLock struct {
  mutex *SpinLock
  owner int
  count int
}

func (re *ReSpinLock) Association(m *SpinLock) {
  re.mutex = m
}

func (re *ReSpinLock) TryLock() bool {
  me := GetCurrentGoroutineId()
  if re.owner == me {
    re.count++
    return true
  }

  return re.mutex.Trylock()
}

func (re *ReSpinLock) Lock() {
  me := GetCurrentGoroutineId()
  if re.owner == me {
    re.count++
    return
  }

  re.mutex.Lock()
}

func (re *ReSpinLock) Unlock() {
  Assert(re.owner == GetCurrentGoroutineId(), "illegalMonitorStateError")
  if re.count > 0 {
    re.count--
  } else {
    re.mutex.Unlock()
  }
}
