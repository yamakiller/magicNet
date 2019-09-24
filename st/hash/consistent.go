package hash

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

func getHash(str []byte) uint32 {
	p := uint32(16777619)
	hash := uint32(2166136261)
	for i := 0; i < len(str); i++ {
		hash = (hash ^ uint32(str[i])) * p
	}
	hash += hash << 13
	hash ^= hash >> 7
	hash += hash << 3
	hash ^= hash >> 17
	hash += hash << 5
	return hash
}

//New Create a hash consistent loader
func New(replicas int) *Map {
	if replicas <= 0 {
		replicas = 20
	}
	return &Map{replicas: replicas,
		circle: make(map[uint32]interface{})}
}

//ErrEmptyCircle 返回一个空
var ErrEmptyCircle = errors.New("empty circle")

type uInt32Slice []uint32

func (s uInt32Slice) Len() int {
	return len(s)
}

func (s uInt32Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s uInt32Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Map Hash consistency load balancing
type Map struct {
	replicas     int
	sortedHashes uInt32Slice
	circle       map[uint32]interface{}
	sync.RWMutex
}

//UnAdd Join an object, not locked
func (m *Map) UnAdd(key string, v interface{}) {
	for i := 0; i < m.replicas; i++ {
		m.circle[getHash([]byte(strconv.Itoa(i)+key))] = v
	}
	m.updateSortedHashes()
}

//UnRemove Delete an object, not locked
func (m *Map) UnRemove(key string) {
	for i := 0; i < m.replicas; i++ {
		delete(m.circle, getHash([]byte(strconv.Itoa(i)+key)))
	}
	m.updateSortedHashes()
}

//UnGet Return an object, not locked
func (m *Map) UnGet(name string) (interface{}, error) {
	if len(m.circle) == 0 {
		return "", ErrEmptyCircle
	}

	key := getHash([]byte(name))
	i := m.sreach(key)
	return m.circle[m.sortedHashes[i]], nil
}

func (m *Map) sreach(key uint32) (i int) {
	f := func(x int) bool {
		return m.sortedHashes[x] > key
	}
	i = sort.Search(len(m.sortedHashes), f)
	if i >= len(m.sortedHashes) {
		i = 0
	}
	return
}

func (m *Map) updateSortedHashes() {
	hashes := m.sortedHashes[:0]
	if cap(m.sortedHashes)/(m.replicas*4) > len(m.circle) {
		hashes = nil
	}
	for k := range m.circle {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	m.sortedHashes = hashes
}
