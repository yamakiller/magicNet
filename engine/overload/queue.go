package overload

import (
	"sync"

	"github.com/yamakiller/magicLibs/util"
)

//Queue desc
//@struct Queue desc : Automatic extension queue
//@member (int) queue cap size
//@member (int) queue head pos
//@member (int) queue tail pos
//@member (int) queue overload of number
//@member (int) queue overlaod threshold
//@member ([]interface{}) queue pools
//@member (sync.Mutex)
type Queue struct {
	cap  int
	head int
	tail int

	overload          int
	overloadThreshold int

	buffer []interface{}
	lock   sync.Mutex
}

//NewQueue desc
//@method NewQueue desc: Create a new Automatic extension queue
//@param  (int) initial size
//@return (*Queue)
func NewQueue(initialSize int) *Queue {
	return &Queue{cap: initialSize,
		head:              0,
		tail:              0,
		overload:          0,
		overloadThreshold: initialSize * 2,
		buffer:            make([]interface{}, initialSize)}
}

//Push desc
//@method Push desc: Insert an object
//@param (interface{}) item
func (q *Queue) Push(item interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.unpush(item)
}

//Pop desc
//@method Pop desc: Take an object, If empty return nil
//@return (interface{}) return object
//@return (bool)
func (q *Queue) Pop() (interface{}, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.unpop()
}

//Overload desc
//@method Overload desc: Detecting queues exceeding the limit [mainly used for warning records]
//@return (int)
func (q *Queue) Overload() int {
	if q.overload != 0 {
		overload := q.overload
		q.overload = 0
		return overload
	}
	return 0
}

//Length desc
//@method Length desc: Length of the queue
//@return (int) length
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
	}

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
