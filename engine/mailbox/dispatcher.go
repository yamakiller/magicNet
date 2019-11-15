package mailbox

import "github.com/yamakiller/magicLibs/coroutine"

//Dispatcher desc
//@interface Dispatcher desc: Publisher interface
type Dispatcher interface {
	Schedule(fn func([]interface{}))
	Throughput() int
}

type goroutineDispatcher int

//Schedule desc
//@method Schedule desc
//@param (func([]interface{})) Running function
func (d goroutineDispatcher) Schedule(fn func([]interface{})) {
	coroutine.Instance().Go(fn)
}

func (d goroutineDispatcher) Throughput() int {
	return int(d)
}

//NewGoroutineDispatcher desc
//@method NewGoroutineDispatcher desc: Create a distributor with a coroutine
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

//NewSynchronizedDispatcher desc
//@method NetSynchronizedDispatcher desc: Create a new synchronous distributor
func NewSynchronizedDispatcher(throughput int) Dispatcher {
	return synchronizedDispatcher(throughput)
}
