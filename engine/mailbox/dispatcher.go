package mailbox

// Dispatcher : 分发器接口
type Dispatcher interface {
	Schedule(fn func())
	Throughput() int
}

type goroutineDispatcher int

func (d goroutineDispatcher) Schedule(fn func()) {
	go fn()
}

func (d goroutineDispatcher) Throughput() int {
	return int(d)
}

// NewGoroutineDispatcher ： 创建一个带协程的分发器
func NewGoroutineDispatcher(throughput int) Dispatcher {
	return goroutineDispatcher(throughput)
}

type synchronizedDispatcher int

func (synchronizedDispatcher) Schedule(fn func()) {
	fn()
}

func (d synchronizedDispatcher) Throughput() int {
	return int(d)
}

// NewSynchronizedDispatcher : 新建一个同步的分发器
func NewSynchronizedDispatcher(throughput int) Dispatcher {
	return synchronizedDispatcher(throughput)
}
