package util

import (
  "sync"
)

type ReMutex struct {
  mutex *sync.Mutex
  owner int
  count int
}

func (re *ReMutex) Association(m *sync.Mutex) {
  re.mutex = m
}

func (re *ReMutex) Lock() {
  me := GetCurrentGoroutineId()
  if re.owner == me {
    re.count++
    return
  }

  re.mutex.Lock()
}

func (re *ReMutex) UnLock() {
  Assert(re.owner == GetCurrentGoroutineId(), "illegalMonitorStateError")
  if re.count > 0 {
    re.count--
  } else {
    re.mutex.Unlock()
  }
}
