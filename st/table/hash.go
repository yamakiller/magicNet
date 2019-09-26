package table

import (
	"errors"

	"github.com/yamakiller/magicNet/st/comparator"
)

var (
	//ErrHashTableFulled Table is full
	ErrHashTableFulled = errors.New("Table is full")
)

//HashTable Hash allocation table
type HashTable struct {
	Mask   uint32
	Max    uint32
	Comp   comparator.Comparator
	seqID  uint32
	sz     int
	arrays []interface{}
}

// Init Initialize the hashtable
func (ht *HashTable) Init() {
	ht.arrays = make([]interface{}, ht.Max)
	ht.seqID = 1
}

//Size returns the hashtable is number
func (ht *HashTable) Size() int {
	return ht.sz
}

//Push Insert an value
func (ht *HashTable) Push(v interface{}) (uint32, error) {
	var i uint32
	for i = 0; i < ht.Max; i++ {
		key := ((i + ht.seqID) & ht.Mask)
		hash := key & (ht.Max - 1)
		if ht.arrays[hash] == nil {
			ht.seqID = key + 1
			ht.arrays[hash] = v
			ht.sz++
			return uint32(key), nil
		}
	}

	return 0, ErrHashTableFulled
}

//Get returns the one elements from the hashtable
func (ht *HashTable) Get(key uint32) interface{} {
	hash := key & uint32(ht.Max-1)
	if ht.arrays[hash] != nil && ht.Comp(ht.arrays[hash], key) == 0 {
		return ht.arrays[hash]
	}
	return nil
}

//GetValues returns the elements of all from hashtable
func (ht *HashTable) GetValues() []interface{} {
	if ht.sz == 0 {
		return nil
	}

	i := 0
	result := make([]interface{}, ht.sz)
	for _, v := range ht.arrays {
		if v == nil {
			continue
		}
		result[i] = ht.arrays[i]
		i++
	}

	return result
}

//Remove removes one elements in the hashtable
func (ht *HashTable) Remove(key uint32) bool {
	hash := uint32(key) & uint32(ht.Max-1)
	if ht.arrays[hash] != nil && ht.Comp(ht.arrays[hash], key) == 0 {
		ht.arrays[hash] = nil
		ht.sz--
		return true
	}

	return false
}
