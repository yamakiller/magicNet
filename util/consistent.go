package util

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uintArray []uint32

func (x uintArray) Len() int {
	return len(x)
}

func (x uintArray) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uintArray) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// Consistent Hash consistency
type Consistent struct {
	NumberOfReplicas int
	sortedHashes     uintArray
	//
	circle map[uint32]interface{}
	size   int
	sync.RWMutex
}

//NewConsistent xxx
func NewConsistent(n int) *Consistent {
	return &Consistent{NumberOfReplicas: n, circle: make(map[uint32]interface{})}
}

//Push inserts a sring element in the consistent hash.
func (slf *Consistent) Push(e string, v interface{}) {
	slf.Lock()
	defer slf.Unlock()
	slf.push(e, v)
}

// Erase removes an element in the consistent hash.
func (slf *Consistent) Erase(e string) {
	slf.Lock()
	defer slf.Unlock()
	slf.erase(e)
}

// Get returns an element close to where name hashes to in the circle
func (slf *Consistent) Get(name string) (interface{}, error) {
	slf.RLock()
	defer slf.RUnlock()

	if len(slf.circle) == 0 {
		return nil, errors.New("empty consistent")
	}
	key := slf.hashCalc(name)
	i := slf.search(key)
	return slf.circle[slf.sortedHashes[i]], nil
}

// Sreach return an element to where f returns 0 to in the circle
func (slf *Consistent) Sreach(key interface{}, f func(key interface{}, val interface{}) int) interface{} {
	slf.RLock()
	defer slf.RUnlock()

	for _, v := range slf.circle {
		if f(key, v) == 0 {
			return v
		}
	}
	return nil
}

// Range Traverse access to all elements
func (slf *Consistent) Range(f func(val interface{})) {
	slf.RLock()
	defer slf.RUnlock()

	for _, v := range slf.circle {
		f(v)
	}
}

//Size returns memory number to in the circle
func (slf *Consistent) Size() int {
	return slf.size
}

func (slf *Consistent) push(e string, v interface{}) {
	for i := 0; i < slf.NumberOfReplicas; i++ {
		slf.circle[slf.hashCalc(slf.genKey(e, i))] = v
	}

	slf.updateSortedHashes()
	slf.size++
}

func (slf *Consistent) erase(e string) {
	for i := 0; i < slf.NumberOfReplicas; i++ {
		delete(slf.circle, slf.hashCalc(slf.genKey(e, i)))
	}
	slf.updateSortedHashes()
	slf.size--
}

func (slf *Consistent) search(key uint32) (i int) {
	f := func(x int) bool {
		return slf.sortedHashes[x] > key
	}

	i = sort.Search(len(slf.sortedHashes), f)
	if i >= len(slf.sortedHashes) {
		i = 0
	}
	return i
}

// generates a string key for an element with an index
func (slf *Consistent) genKey(s string, idx int) string {
	return strconv.Itoa(idx) + s
}

func (slf *Consistent) hashCalc(s string) uint32 {
	if len(s) < 64 {
		var scratch [64]byte
		copy(scratch[:], s)
		return crc32.ChecksumIEEE(scratch[:len(s)])
	}
	return crc32.ChecksumIEEE([]byte(s))
}

func (slf *Consistent) updateSortedHashes() {
	hashes := slf.sortedHashes[:0]
	if cap(slf.sortedHashes)/(slf.NumberOfReplicas*4) > len(slf.circle) {
		hashes = nil
	}
	for k := range slf.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	slf.sortedHashes = hashes
}
