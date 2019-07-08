package actor

import (
	"magicNet/engine/util"
	"sync"
	"sync/atomic"
)

// Registry ：注册表
type Registry struct {
	localAddress   uint32
	localSequence  uint32
	localItem      []*PID
	localItemMutex sync.RWMutex
}

const (
	registerDefaultSize = 32
)

// GlobalRegistry : Actor 全局注册表
var GlobalRegistry = &Registry{
	localItem: make([]*PID, registerDefaultSize),
}

// SetLocalAddress : 设置本地服务器地址
func (r *Registry) SetLocalAddress(addr uint32) {
	r.localAddress = addr
}

// GetLocalAddress : 获取本地服务器地址信息
func (r *Registry) GetLocalAddress() uint32 {
	return r.localAddress
}

// Register : 注册一个Actor并生成PID
func (r *Registry) Register(pid *PID) bool {
	r.localItemMutex.Lock()
	for {
		var i uint32
		currentNum := uint32(len(r.localItem))
		for i = 0; i < currentNum; i++ {
			key := ((i + r.localSequence) & pidMask)
			hash := key & (currentNum - 1)
			if r.localItem[hash] == nil {
				r.localItem[hash] = pid
				r.localSequence = key + 1
				r.localItemMutex.Unlock()
				pid.ID = (key | (r.localAddress << pidKeyBit))
				return true
			}
		}

		newNum := (currentNum * 2)
		util.Assert(newNum <= pidMax, "actor number overflow")
		newItem := make([]*PID, newNum)

		for i = 0; i < currentNum; i++ {
			if newItem[i] == nil {
				continue
			}

			hash := (newItem[i].Key() & (newNum - 1))
			if hash == i {
				continue
			}
			newItem[hash] = r.localItem[i]
		}

		r.localItem = newItem
	}
}

// UnRegister : 注销 PID
func (r *Registry) UnRegister(pid *PID) bool {
	r.localItemMutex.Lock()
	defer r.localItemMutex.Unlock()
	hash := pid.Key() & uint32(len(r.localItem)-1)
	if r.localItem[hash] != nil && r.localItem[hash].Equal(pid) {
		ref := r.localItem[hash].p
		if l, ok := (*ref).(*AtrProcess); ok {
			atomic.StoreInt32(&l.death, 1)
		}
		r.localItem[hash] = nil
		return true
	}
	return false
}

// Get : 获取PID的处理对象
func (r *Registry) Get(pid *PID) (Process, bool) {
	r.localItemMutex.RLock()
	defer r.localItemMutex.RUnlock()

	hash := pid.Key() & uint32(len(r.localItem)-1)
	if r.localItem[hash] != nil && r.localItem[hash].Equal(pid) {
		return *r.localItem[hash].p, true
	}
	return deathLetter, false
}
