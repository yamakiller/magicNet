package overload

import (
	"magicNet/engine/util"
	"sync"
)

// Queue : 无上限队列
type Queue struct {
	cap  int
	head int
	tail int

	overload          int
	overloadThreshold int

	buffer []interface{}
	lock   sync.Mutex
}

// NewQueue : 新建一个无上限队列
func NewQueue(initialSize int) *Queue {
	return &Queue{cap: initialSize,
		head:              0,
		tail:              0,
		overload:          0,
		overloadThreshold: initialSize * 2,
		buffer:            make([]interface{}, initialSize)}
}

// Push : 插入一个对象
func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.unpush(item)
}

// Pop : 取出一个对象, If empty return nil
func (q *Queue) Pop() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.unpop()
}

// Overload : 检测队列超出限制的情况[主要用于警告记录]
func (q *Queue) Overload() int {
	if q.overload != 0 {
		overload := q.overload
		q.overload = 0
		return overload
	}
	return 0
}

// Length : 队列的长度
func (q *Queue) Length() int {
	var (
		head int
		tail int
		cap  int
	)
	q.lock.Lock()
	head = q.head
	tail = q.tail
	cap = q.cap
	q.lock.Unlock()

	if head <= tail {
		return tail - head
	}
	return tail + cap - head
}

func (q *Queue) unpush(item interface{}) {
	util.AssertEmpty(item, "error push is nil")
	q.buffer[q.tail] = item
	q.tail++
	if q.tail >= q.cap {
		q.tail = 0
	}

	if q.head == q.tail {
		q.expand()
	}
}

func (q *Queue) unpop() (interface{}, bool) {
	var resultSucces bool
	var result interface{}
	if q.head != q.tail {
		resultSucces = true
		result = q.buffer[q.head]
		q.buffer[q.head] = nil
		q.head++
		if q.head >= q.cap {
			q.head = 0
		}

		length := q.tail - q.tail
		if length < 0 {
			length += q.cap
		}
		for length > q.overloadThreshold {
			q.overload = length
			q.overloadThreshold *= 2
		}
	} /*
		 ! 这里是否需要这样呢？
		else {
		  q.overloadThreshold = q.cap
		}*/
	return result, resultSucces
}

func (q *Queue) expand() {
	newBuff := make([]interface{}, q.cap*2)
	for i := 0; i < q.cap; i++ {
		newBuff[i] = q.buffer[(q.head+i)%q.cap]
	}

	q.head = 0
	q.tail = q.cap
	q.cap *= 2

	q.buffer = newBuff
}
