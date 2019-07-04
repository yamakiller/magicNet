package timer

import (
	"container/list"
	"magicNet/engine/monitor"
	"sync"
	"time"
)

const (
	timeNearShift  = 8
	timeNear       = (1 << timeNearShift) // 256
	timeLevelShift = 6
	timeLevel      = (1 << timeLevelShift) // 64
	timeNearMask   = (timeNear - 1)
	timeLevelMask  = (timeLevel - 1)
)

// TimeGoFunction : 定时器执行函数
type TimeGoFunction func(parm interface{})

// Timer : 定时器
type Timer struct {
	near         [timeNear]*list.List
	t            [4][timeLevel]*list.List
	m            sync.Mutex
	timeDur      uint32
	startTime    uint32
	current      uint64
	currentPoint uint64
	shutdown     chan struct{}
}

type Node struct {
	expire uint32
	param  interface{}
	f      TimeGoFunction
}

var instTime *Timer

// StartService  : 初始化定时器对象
func StartService() {
	instTime = &Timer{}
	instTime.timeDur = 0
	instTime.startTime = 0
	instTime.current = 0
	instTime.currentPoint = 0
	instTime.shutdown = make(chan struct{})
	var i, j int
	for i = 0; i < timeNear; i++ {
		instTime.near[i] = list.New()
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < timeLevel; j++ {
			instTime.t[i][j] = list.New()
		}
	}

	instTime.start()
}

// StopService : 关闭定时器
func StopService() {
	instTime.stop()
}

// TimeOut : 设置定时事件
func TimeOut(tm int, f TimeGoFunction, param interface{}) int {
	if tm <= 0 {
		f(param)
	} else {
		instTime.putTime(tm, f, param)
	}
	return 0
}

// StartTime : 获取定时器开始时间
func StartTime() uint32 {
	return instTime.getStartTime()
}

// Now : 获取当前时间
func Now() uint64 {
	return instTime.getNow()
}

// Start : 启动定时器
func (T *Timer) start() {
	tnow := time.Now()
	T.startTime = uint32(tnow.Unix())
	T.current = uint64(tnow.UnixNano() / 10000000)
	T.currentPoint = T.getTime()
	monitor.WaitInc()
	go func(t *Timer) {
		tick := time.NewTicker(time.Nanosecond * 1000)
		defer tick.Stop()
		defer monitor.WaitDec()
		for {
			select {
			case <-tick.C:
				t.update()
			case <-t.shutdown:
				return
			}
		}
	}(T)
}

func (T *Timer) stop() {
	close(T.shutdown)
}

func (T *Timer) getStartTime() uint32 {
	return T.startTime
}

func (T *Timer) getNow() uint64 {
	return T.current
}

func (T *Timer) update() {
	T.m.Lock()

	T.execute()

	T.shift()

	T.execute()

	T.m.Unlock()
}

func dispatchList(front *list.Element) {
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		go node.f(node.param)
	}
}

func (T *Timer) execute() {
	idx := T.timeDur & timeNearMask
	for {
		if T.near[idx].Len() <= 0 {
			break
		}

		front := T.near[idx].Front()
		T.near[idx].Init()
		T.m.Unlock()
		dispatchList(front)
		T.m.Lock()
	}
}

func (T *Timer) shift() {
	var mask uint32 = timeNear
	T.timeDur++
	ct := T.timeDur
	if ct == 0 {
		T.moveList(3, 0)
	} else {
		timeDur := ct >> timeNearShift
		i := 0
		for (ct & (mask - 1)) == 0 {

			idx := int(timeDur & timeLevelMask)
			if idx != 0 {
				T.moveList(i, idx)
				break
			}

			mask <<= timeLevelShift
			timeDur >>= timeLevelShift
			i++

		}
	}
}

func (T *Timer) moveList(level, idx int) {
	front := T.t[level][idx].Front()
	T.t[level][idx].Init()
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		T.addNode(node)
	}
}

func (T *Timer) addNode(n *Node) {
	expire := n.expire
	currentTime := T.timeDur

	if (expire | timeNearMask) == (currentTime | timeNearMask) {
		T.near[expire&timeNearMask].PushBack(n)
	} else {
		var i uint32
		var mask uint32 = timeNear << timeLevelShift
		for i = 0; i < 3; i++ {
			if (expire | (mask - 1)) == (currentTime | (mask - 1)) {
				break
			}
			mask <<= timeLevelShift
		}

		T.t[i][(expire>>(timeNearShift+i*timeLevelShift))&timeLevelMask].PushBack(n)
	}
}

func (T *Timer) putTime(tm int, f TimeGoFunction, parm interface{}) {

	n := &Node{0, parm, f}
	T.m.Lock()
	defer T.m.Unlock()
	n.expire = uint32(tm) + T.timeDur
	T.addNode(n)
}

func (T *Timer) getTime() uint64 {
	return uint64(time.Now().UnixNano() / 10000000)
}
