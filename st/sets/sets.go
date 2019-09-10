package sets

import (
	"github.com/yamakiller/magicNet/st/containers"
)

// Set interface that all sets
type Set interface {
	Push(es ...interface{})
	PushAll(st *Set)
	Retain(eds ...interface{})
	RetainAll(st *Set)
	Erase(es ...interface{})
	EraseAll(st *Set)
	Contains(es ...interface{})

	containers.Container
}
