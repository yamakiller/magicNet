package containers

// Container is base interface that all data structures
type Container interface {
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}
