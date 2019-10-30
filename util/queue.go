package util

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	val  interface{}
}

// Queue : Simple queue
type Queue struct {
	head, tail *node
}

// NewQueue : Create a queue object
func NewQueue() *Queue {
	q := &Queue{}
	stub := &node{}
	q.head = stub
	q.tail = stub
	return q
}

// Push : Insert an Object into the queue
func (slf *Queue) Push(t interface{}) {
	n := new(node)
	n.val = t
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&slf.head)), unsafe.Pointer(n)))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&prev.next)), unsafe.Pointer(n))
}

// Pop : An object pops up in the re-queue
func (slf *Queue) Pop() interface{} {
	tail := slf.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	if next != nil {
		slf.tail = next
		v := next.val
		next.val = nil
		return v
	}
	return nil
}

// IsEmpty : Whether the queue is empty
func (slf *Queue) IsEmpty() bool {
	tail := slf.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}
