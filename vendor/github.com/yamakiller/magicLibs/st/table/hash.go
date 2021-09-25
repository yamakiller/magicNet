package table

import (
	"errors"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/st/comparator"
)

var (
	//ErrHashTableFulled Table is full
	ErrHashTableFulled = errors.New("Table is full")
)

//HashTable Hash allocation table
type HashTable struct {
	Mask    uint32
	Max     uint32
	Comp    comparator.Comparator
	_seqID  uint32
	_sz     int
	_arrays []interface{}
}

//Initial 初始化Hash表
func (slf *HashTable) Initial() {
	slf._arrays = make([]interface{}, slf.Max)
	slf._seqID = 1
	if !mmath.IsPower(int(slf.Max)) {
		panic("parameter Max must be a power of two")
	}
}

//Size 返回元素个数
func (slf *HashTable) Size() int {
	return slf._sz
}

//Push 插入一个元素并分配一个ID
func (slf *HashTable) Push(v interface{}) (uint32, error) {
	var i uint32
	for i = 0; i < slf.Max; i++ {
		key := ((i + slf._seqID) & slf.Mask)
		if key == 0 {
			key = 1
		}
		hash := key & (slf.Max - 1)
		if slf._arrays[hash] == nil {
			slf._seqID = key + 1
			slf._arrays[hash] = v
			slf._sz++
			return uint32(key), nil
		}
	}

	return 0, ErrHashTableFulled
}

//Get 返回一个元素
func (slf *HashTable) Get(key uint32) interface{} {
	hash := key & uint32(slf.Max-1)
	if slf._arrays[hash] != nil && slf.Comp(slf._arrays[hash], key) == 0 {
		return slf._arrays[hash]
	}
	return nil
}

//GetValues 返回所有元素
func (slf *HashTable) GetValues() []interface{} {
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
func (slf *HashTable) Remove(key uint32) bool {
	hash := uint32(key) & uint32(slf.Max-1)
	if slf._arrays[hash] != nil && slf.Comp(slf._arrays[hash], key) == 0 {
		slf._arrays[hash] = nil
		slf._sz--
		return true
	}

	return false
}
