package mailbox

import "github.com/yamakiller/magicLibs/coroutine"

//Dispatcher doc
//@Interface Dispatcher @Summary Publisher interface
type Dispatcher interface {
	Schedule(fn func([]interface{}))
	Throughput() int
}

type goroutineDispatcher int

//Schedule doc
//@Method Schedule doc
//@Param (func([]interface{})) Running function
func (d goroutineDispatcher) Schedule(fn func([]interface{})) {
	coroutine.Instance().Go(fn)
}

func (d goroutineDispatcher) Throughput() int {
	return int(d)
}

//NewGoroutineDispatcher doc
//@Method NewGoroutineDispatcher @Summary Create a distributor with a coroutine
func NewGoroutineDispatcher(throughput int) Dispatcher {
	return goroutineDispatcher(throughput)
}

type synchronizedDispatcher int

func (synchronizedDispatcher) Schedule(fn func([]interface{})) {
	fn(nil)
}

func (d synchronizedDispatcher) Throughput() int {
	return int(d)
}

//NewSynchronizedDispatcher doc
//@Method NetSynchronizedDispatcher @Summary Create a new synchronous distributor
func NewSynchronizedDispatcher(throughput int) Dispatcher {
	return synchronizedDispatcher(throughput)
}
