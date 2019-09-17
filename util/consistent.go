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


func NewConsistent(n int) *Consistent {
	return &Consistent{NumberOfReplicas : n, circle: make(map[uint32]interface{})}
}

//Push inserts a sring element in the consistent hash.
func (c *Consistent) Push(e string, v interface{}) {
	c.Lock()
	defer c.Unlock()
	c.push(e, v)
}

// Erase removes an element in the consistent hash.
func (c *Consistent) Erase(e string) {
	c.Lock()
	defer c.Unlock()
	c.erase(e)
}

// Get returns an element close to where name hashes to in the circle
func (c *Consistent) Get(name string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.circle) == 0 {
		return nil, errors.New("empty consistent")
	}
	key := c.hashCalc(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}

// Sreach return an element to where f returns 0 to in the circle
func (c *Consistent) Sreach(key interface{}, f func(key interface{}, val interface{}) int) interface{} {
	c.RLock()
	defer c.RUnlock()

	for _, v := range c.circle {
		if f(key, v) == 0 {
			return v
		}
	}
	return nil
}

// Range Traverse access to all elements
func (c *Consistent) Range(f func(val interface{})) {
	c.RLock()
	defer c.RUnlock()

	for _, v := range c.circle {
		f(v)
	}
}

//Size returns memory number to in the circle
func (c *Consistent) Size() int {
	return c.size
}

func (c *Consistent) push(e string, v interface{}) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		c.circle[c.hashCalc(c.genKey(e, i))] = v
	}

	c.updateSortedHashes()
	c.size++
}

func (c *Consistent) erase(e string) {
	for i := 0; i < c.NumberOfReplicas; i++ {
		delete(c.circle, c.hashCalc(c.genKey(e, i)))
	}
	c.updateSortedHashes()
	c.size--
}

func (c *Consistent) search(key uint32) (i int) {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}

	i = sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

// generates a string key for an element with an index
func (c *Consistent) genKey(s string, idx int) string {
	return strconv.Itoa(idx) + s
}

func (c *Consistent) hashCalc(s string) uint32 {
	if len(s) < 64 {
		var scratch [64]byte
		copy(scratch[:], s)
		return crc32.ChecksumIEEE(scratch[:len(s)])
	}
	return crc32.ChecksumIEEE([]byte(s))
}

func (c *Consistent) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.NumberOfReplicas*4) > len(c.circle) {
		hashes = nil
	}
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedHashes = hashes
}
