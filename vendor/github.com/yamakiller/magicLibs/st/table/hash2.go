package table

import (
	"sync"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/st/comparator"
)

//HashTable2 自动放大的Hash表
type HashTable2 struct {
	Mask    uint32
	Max     uint32 //2的幂
	Comp    comparator.Comparator
	GetKey  func(a interface{}) uint32
	_seqID  uint32
	_sz     int
	_cap    uint32
	_arrays []interface{}
	_sync   sync.Mutex
}

//Initial 初始化Hash表
func (slf *HashTable2) Initial() {
	slf._cap = 16
	slf._arrays = make([]interface{}, slf._cap)
	slf._seqID = 1

	if !mmath.IsPower(int(slf.Max)) {
		panic("parameter Max must be a power of two")
	}
}

//Size 返回元素数量
func (slf *HashTable2) Size() int {
	return slf._sz
}

//Push 插入一个元素，并分配一个ID
func (slf *HashTable2) Push(v interface{}) (uint32, error) {
	var i uint32

	slf._sync.Lock()
	for {

		for i = 0; i < slf._cap; i++ {
			key := ((i + slf._seqID) & slf.Mask)
			if key == 0 {
				key = 1
			}
			hash := key & (slf._cap - 1)
			if slf._arrays[hash] == nil {
				slf._seqID = key + 1
				slf._arrays[hash] = v
				slf._sz++
				slf._sync.Unlock()
				return uint32(key), nil
			}
		}

		newCap := slf._cap * 2
		if newCap > slf.Max {
			newCap = slf.Max
		}
		if newCap == slf._cap {
			slf._sync.Unlock()
			return 0, ErrHashTableFulled
		}

		slf._arrays = append(slf._arrays, make([]interface{}, newCap-slf._cap)...)
		for i = 0; i < slf._cap; i++ {
			if slf._arrays[i] == nil {
				continue
			}

			hash := slf.GetKey(slf._arrays[i]) & uint32(newCap-1)
			if hash == i {
				continue
			}

			tmp := slf._arrays[i]
			slf._arrays[i] = nil
			slf._arrays[hash] = tmp
		}
		slf._cap = newCap
	}
}

//Get 获取一个元素
func (slf *HashTable2) Get(key uint32) interface{} {
	slf._sync.Lock()
	defer slf._sync.Unlock()

	hash := key & uint32(slf._cap-1)
	if slf._arrays[hash] != nil && slf.Comp(slf._arrays[hash], key) == 0 {
		return slf._arrays[hash]
	}
	return nil
}

//GetValues 返回所有元素
func (slf *HashTable2) GetValues() []interface{} {
	slf._sync.Lock()
	defer slf._sync.Unlock()

	if slf._sz == 0 {
		return nil
	}

	i := 0
	result := make([]interface{}, slf._sz)
	for _, v := range slf._arrays {
		if v == nil {
			continue
		}
		result[i] = v
		i++
	}

	return result
}

//Remove 删除一个元素
func (slf *HashTable2) Remove(key uint32) bool {
	slf._sync.Lock()
	defer slf._sync.Unlock()

	hash := uint32(key) & uint32(slf._cap-1)
	if slf._arrays[hash] != nil && slf.Comp(slf._arrays[hash], key) == 0 {
		slf._arrays[hash] = nil
		slf._sz--
		return true
	}

	return false
}
