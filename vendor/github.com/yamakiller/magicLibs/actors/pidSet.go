package actors

import (
	"errors"
	"math"
	"time"

	"github.com/yamakiller/magicLibs/mutex"
)

//DPIDSet PID 集
type PIDSet struct {
	_pids []PID
	_seq  int
	_sz   int
	_syn  mutex.SpinLock
}

//Next 生成新的PID
func (slf *PIDSet) Next() (*PID, error) {
	slf.lock()
	defer slf.unlock()

	for i := 0; i < math.MaxUint16; i++ {
		key := ((i + slf._seq) & math.MaxUint16)
		hash := key & (math.MaxUint16 - 1)
		if slf._pids[hash].ID == 0 {
			slf._seq = key + 1
			slf._pids[hash].ID = uint32(key)
			slf._sz++
			return &slf._pids[hash], nil
		}
	}

	return &PID{}, errors.New("full error")
}

//Remove 移出ID
func (slf *PIDSet) Remove(pid *PID) error {
	slf.lock()
	defer slf.unlock()

	hash := uint32(pid.ID) & uint32(math.MaxUint16-1)
	if slf._pids[hash].ID != 0 && slf._pids[hash].ID == pid.ID {
		slf._pids[hash].ID = 0
		slf._sz--
		return nil
	}

	return errors.New("not fount")
}

//Values 获取所有的PID
func (slf *PIDSet) Values() []*PID {
	slf.lock()
	defer slf.unlock()

	icur := 0
	result := make([]*PID, slf._sz)
	for i := 0; i < math.MaxUint16; i++ {
		if slf._pids[i].ID > 0 {
			result[icur] = &slf._pids[i]
			icur++
			if icur >= slf._sz {
				break
			}
		}
	}

	return result
}

func (slf *PIDSet) lock() {
	try := 0
	timeDelay := time.Millisecond
	for {
		if !slf._syn.Trylock() {
			try++
			if try > 6 {
				time.Sleep(timeDelay)
				timeDelay = timeDelay * 2
				if max := 100 * time.Microsecond; timeDelay > max {
					timeDelay = max
				}
			}
			continue
		}
		break
	}
}

func (slf *PIDSet) unlock() {
	slf._syn.Unlock()
}
