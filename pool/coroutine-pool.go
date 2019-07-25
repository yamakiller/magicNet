package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type coState int32
type coFunc func([]interface{})

const (
	coIdle    = int32(0)
	coRun     = int32(1)
	coClosing = int32(2)
	coDeath   = int32(3)
)

var (
	// ErrPoolStoped : 协程池已关闭
	ErrPoolStoped = errors.New("coroutine pool stoped")
)

type task struct {
	cb     coFunc
	params []interface{}
}

func (co *coObject) run(cop *CoroutinePool) {
	go func() {
		defer cop.wait.Done()
		for {
			for {
				select {
				case <-co.q:
					co.state = coDeath
					return
				case t := <-cop.taskQueue:
					atomic.CompareAndSwapInt32(&co.state, coIdle, coRun)
					t.cb(t.params)
					t.params = nil
					atomic.CompareAndSwapInt32(&co.state, coRun, coIdle)
				}
			}
		}
	}()
}

type coObject struct {
	state int32
	last  time.Duration
	q     chan int
	id    int
}

/*var (
	taskLimit int
	maxCoNum  int
	minCoNum  int

	taskQueue chan task
	sep       int
	cos       []coObject
	wait      sync.WaitGroup
	quit      chan int
)*/

// CoroutinePool : 协程池
type CoroutinePool struct {
	TaskLimit int
	MaxNum    int
	MinNum    int

	taskQueue chan task
	runing    int32
	curr      int32
	seq       int32
	cos       []coObject
	wait      sync.WaitGroup
	quit      chan int
}

func (co *CoroutinePool) scheduer() {
	defer co.wait.Done()
	for {
		select {
		case <-co.quit:
			goto scheduer_end
		default:
		}

		now := time.Now().Unix()
		for _, v := range co.cos {
			if v.state == coRun ||
				v.state == coDeath ||
				v.state == coClosing {
				continue
			}

			if v.state == coIdle && ((time.Duration(now) - v.last) > time.Second*60) {
				if atomic.CompareAndSwapInt32(&v.state, int32(coIdle), int32(coClosing)) {
					close(v.q)
				}
			}
		}

		time.Sleep(time.Millisecond * 1000)
	}
scheduer_end:
	for i := 0; i < co.MaxNum; i++ {
		if co.cos[i].state != coDeath && co.cos[i].state != coClosing {
			atomic.StoreInt32(&co.cos[i].state, coClosing)
			close(co.cos[i].q)
		}
	}
}

// Start : 启动协程池
func (co *CoroutinePool) Start() {
	if co.MaxNum == 0 {
		co.MaxNum = 65535
	}
	atomic.StoreInt32(&co.seq, 1)

	now := time.Now().Unix()
	co.taskQueue = make(chan task, co.TaskLimit)
	co.cos = make([]coObject, co.MaxNum)
	for k, v := range co.cos {
		v.id = k + 1
		v.state = coDeath
		v.last = time.Duration(now)
	}

	if co.MaxNum > 0 && co.MaxNum > co.MinNum {
		co.wait.Add(1)
		go co.scheduer()
	}

	for i := 0; i < co.MinNum; i++ {
		co.cos[i].state = coIdle
		co.startOne(i)
	}
}

func (co *CoroutinePool) startOne(idx int) {
	co.wait.Add(1)
	atomic.AddInt32(&co.seq, 1)
	atomic.AddInt32(&co.curr, 1)

	co.cos[idx].run(co)
}

// StopPool : 关闭协程池
func (co *CoroutinePool) StopPool() {
	close(co.quit)
	co.wait.Wait()
}

//Go 运行任务
func (co *CoroutinePool) Go(f func(params []interface{}), params ...interface{}) error {

	select {
	case <-co.quit:
		return ErrPoolStoped
	default:
	}
	runing := atomic.LoadInt32(&co.runing)
	curr := atomic.LoadInt32(&co.curr)

	if runing >= curr && curr < int32(co.MaxNum) {
		for i := 0; i < co.MaxNum; i++ {
			hash := ((int32(i) + co.seq) % int32(co.MaxNum))
			if atomic.CompareAndSwapInt32(&co.cos[hash].state, coDeath, coIdle) {
				co.seq = int32(hash) + co.seq
				co.startOne(int(hash))
				break
			}
		}
	}

	select {
	case <-co.quit:
		return ErrPoolStoped
	case co.taskQueue <- task{f, params}:
		return nil
	}
}
