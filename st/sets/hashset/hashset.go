package hashset

import (
	"fmt"
	"strings"
)

// Set holds elements in go`s native map
type Set struct {
	items map[interface{}]struct{}
}

// Push adds the es (one or more) to the set
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.items[it] = struct{}{}
	}
}

// PushAll adds the st(set in element) to the set.
func (s *Set) PushAll(st *Set) {
	for _, it := range st.items {
		s.items[it] = struct{}{}
	}
}

// Retain retain the es (one or more) to the set.
func (s *Set) Retain(eds ...interface{}) {
	vs := make(map[interface{}]struct{})
	for _, it := range eds {
		if v, ok := s.items[it]; ok {
			vs[it] = v
		}
	}

	s.items = vs
}

// RetainAll retain the st(set in element) to the set.
func (s *Set) RetainAll(st *Set) {
	vs := make(map[interface{}]struct{})
	for _, it := range st.items {
		if v, ok := s.items[it]; ok {
			vs[it] = v
		}
	}
	s.items = vs
}

// Erase removes the es (one or more) from the set
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		delete(s.items, it)
	}
}

// EraseAll removes this st(set in element) from the set
func (s *Set) EraseAll(st *Set) {
	for _, it := range st.items {
		delete(s.items, it)
	}
}

// Contains check if es (one or more) are present in the set.
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s.items[it]; !cs {
			return false
		}
	}
	return true
}

// Size returns number of elements within the set.
func (s *Set) Size() int {
	return len(s.items)
}

// Empty returns true if set does not contain any elements.
func (s *Set) Empty() bool {
	return s.Size() == 0
}

// Clear clears all values in the set.
func (s *Set) Clear() {
	s.items = make(map[interface{}]struct{})
}

// Values returns all items in the set.
func (s *Set) Values() []interface{} {
	vs := make([]interface{}, s.Size())
	icnt := 0
	for it := range s.items {
		vs[icnt] = it
		icnt++
	}
	return vs
}

// String returns a string
func (s *Set) String() string {
	str := "HashSet\n"
	items := []string{}
	for k := range s.items {
		items = append(items, fmt.Sprintf("%v", k))
	}
	str += strings.Join(items, ", ")
	return str
}
