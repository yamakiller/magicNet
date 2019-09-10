package treeset

import (
	"fmt"
	"strings"

	rbt "github.com/yamakiller/magicNet/st/trees/redblacktree"
)

// Set holds elements in a red-black tree
type Set struct {
	tree *rbt.Tree
}

// Push adds the es (one or more) to the set
func (s *Set) Push(es ...interface{}) {
	for _, it := range es {
		s.tree.Insert(it, struct{}{})
	}
}

// PushAll adds the st(set in element) to the set.
func (s *Set) PushAll(st *Set) {
	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s.tree.Insert(it.Key(), it.Value())
		}
	}
}

// Retain retain the es (one or more) to the set.
func (s *Set) Retain(eds ...interface{}) {
	var vs []interface{}
	var ic int
	for _, it := range eds {
		if _, ok := s.tree.Get(it); ok {
			vs = append(vs, it)
			ic++
		}
	}

	s.tree.Clear()
	for i := 0; i < ic; i++ {
		s.tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

// RetainAll retain the st(set in element) to the set.
func (s *Set) RetainAll(st *Set) {
	var vs []interface{}
	var ic int

	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			if _, ok := s.tree.Get(it.Key()); ok {
				vs = append(vs, it)
				ic++
			}
		}
	}

	s.tree.Clear()
	for i := 0; i < ic; i++ {
		s.tree.Insert(vs[i], struct{}{})
	}
	vs = nil
}

// Erase removes the es (one or more) from the set
func (s *Set) Erase(es ...interface{}) {
	for _, it := range es {
		s.tree.Erase(it)
	}
}

// EraseAll removes this st(set in element) from the set
func (s *Set) EraseAll(st *Set) {
	it := st.tree.Iterator()
	for i := 0; it.Next(); i++ {
		if it.It() != nil {
			s.tree.Erase(it.Key())
		}
	}
}

// Contains check if es (one or more) are present in the set.
func (s *Set) Contains(es ...interface{}) bool {
	for _, it := range es {
		if _, cs := s.tree.Get(it); !cs {
			return false
		}
	}
	return true
}

// Size returns number of elements within the set.
func (s *Set) Size() int {
	return s.tree.Size()
}

// Empty returns true if set does not contain any elements.
func (s *Set) Empty() bool {
	return s.Size() == 0
}

// Clear clears all values in the set.
func (s *Set) Clear() {
	s.tree.Clear()
}

// Values returns all items in the set.
func (s *Set) Values() []interface{} {
	return s.tree.Keys()
}

// String returns a string
func (s *Set) String() string {
	str := "TreeSet\n"
	items := []string{}
	for _, v := range s.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}
