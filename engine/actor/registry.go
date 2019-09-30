package actor

import (
	"sync"
	"sync/atomic"

	"github.com/yamakiller/magicNet/util"
)

type registryValue struct {
	key uint32
	val interface{}
}

// Registry ：注册表
type Registry struct {
	localAddress   uint32
	localSequence  uint32
	localItem      []registryValue
	localItemMutex sync.RWMutex
}

const (
	registerDefaultSize = 32
)

// globalRegistry : Global registry
var globalRegistry = &Registry{
	localSequence: 1,
	localItem:     make([]registryValue, registerDefaultSize),
}

// SetLocalAddress Set the local server address
func (r *Registry) SetLocalAddress(addr uint32) {
	r.localAddress = addr
}

// GetLocalAddress Get local server address information
func (r *Registry) GetLocalAddress() uint32 {
	return r.localAddress
}

// Register Register an Actor and generate a PID
func (r *Registry) Register(pid *PID, process Process) bool {
	r.localItemMutex.Lock()
	for {
		var i uint32
		currentNum := uint32(len(r.localItem))
		for i = 0; i < currentNum; i++ {
			key := ((i + r.localSequence) & pidMask)
			hash := key & (currentNum - 1)
			if r.localItem[hash].key == 0 {
				r.localItem[hash].key = key
				r.localItem[hash].val = process
				r.localSequence = key + 1
				r.localItemMutex.Unlock()
				pid.ID = (key | (r.localAddress << pidKeyBit))
				return true
			}
		}

		newNum := (currentNum * 2)
		util.Assert(newNum <= pidMax, "actor number overflow")
		newItem := make([]registryValue, newNum)

		for i = 0; i < currentNum; i++ {
			if r.localItem[i].key == 0 {
				continue
			}

			hash := (r.localItem[i].key & (newNum - 1))
			if hash == i {
				continue
			}
			newItem[hash] = r.localItem[i]
		}

		r.localItem = newItem
	}
}

// UnRegister Logout PID
func (r *Registry) UnRegister(pid *PID) bool {
	r.localItemMutex.Lock()
	defer r.localItemMutex.Unlock()
	hash := pid.Key() & uint32(len(r.localItem)-1)
	if r.localItem[hash].key != 0 && r.localItem[hash].key == pid.Key() {
		ref := r.localItem[hash].val
		if l, ok := ref.(*AtrProcess); ok {
			atomic.StoreInt32(&l.death, 1)
		}
		r.localItem[hash].key = 01
		r.localItem[hash].val = nil
		return true
	}
	return false
}

// Get  Return the processing object of the PID
func (r *Registry) Get(pid *PID) (Process, bool) {
	r.localItemMutex.RLock()
	defer r.localItemMutex.RUnlock()

	hash := pid.Key() & uint32(len(r.localItem)-1)
	if r.localItem[hash].key != 0 && r.localItem[hash].key == pid.Key() {
		return r.localItem[hash].val.(Process), true
	}
	return deathLetter, false
}
