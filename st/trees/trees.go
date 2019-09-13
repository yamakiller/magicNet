package trees

import (
	"github.com/emirpasic/gods/containers"
)

// Tree interface that all trees
type Tree interface {
	Insert(k, v interface{})
	Erase(k interface{})
	Get(k interface{}) (interface{}, bool)

	containers.Container
}
