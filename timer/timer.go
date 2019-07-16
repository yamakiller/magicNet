package timer

import (
	"container/list"
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
	w            sync.WaitGroup
	timeDur      uint32
	startTime    uint32
	current      uint64
	currentPoint uint64
	shutdown     chan struct{}
}

// Node : 定时器节点
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
func (tm *Timer) start() {
	tnow := time.Now()
	tm.startTime = uint32(tnow.Unix())
	tm.current = uint64(tnow.UnixNano() / 10000000)
	tm.currentPoint = tm.getTime()

	tm.w.Add(1)
	go func(t *Timer) {
		tick := time.NewTicker(time.Nanosecond * 1000)
		defer tick.Stop()
		defer tm.w.Done()
		for {
			select {
			case <-tick.C:
				t.update()
			case <-t.shutdown:
				return
			}
		}
	}(tm)
}

func (tm *Timer) stop() {
	close(tm.shutdown)
}

func (tm *Timer) getStartTime() uint32 {
	return tm.startTime
}

func (tm *Timer) getNow() uint64 {
	return tm.current
}

func (tm *Timer) update() {
	tm.m.Lock()

	tm.execute()

	tm.shift()

	tm.execute()

	tm.m.Unlock()
}

func dispatchList(front *list.Element) {
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		go node.f(node.param)
	}
}

func (tm *Timer) execute() {
	idx := tm.timeDur & timeNearMask
	for {
		if tm.near[idx].Len() <= 0 {
			break
		}

		front := tm.near[idx].Front()
		tm.near[idx].Init()
		tm.m.Unlock()
		dispatchList(front)
		tm.m.Lock()
	}
}

func (tm *Timer) shift() {
	var mask uint32 = timeNear
	tm.timeDur++
	ct := tm.timeDur
	if ct == 0 {
		tm.moveList(3, 0)
	} else {
		timeDur := ct >> timeNearShift
		i := 0
		for (ct & (mask - 1)) == 0 {

			idx := int(timeDur & timeLevelMask)
			if idx != 0 {
				tm.moveList(i, idx)
				break
			}

			mask <<= timeLevelShift
			timeDur >>= timeLevelShift
			i++

		}
	}
}

func (tm *Timer) moveList(level, idx int) {
	front := tm.t[level][idx].Front()
	tm.t[level][idx].Init()
	for e := front; e != nil; e = e.Next() {
		node := e.Value.(*Node)
		tm.addNode(node)
	}
}

func (tm *Timer) addNode(n *Node) {
	expire := n.expire
	currentTime := tm.timeDur

	if (expire | timeNearMask) == (currentTime | timeNearMask) {
		tm.near[expire&timeNearMask].PushBack(n)
	} else {
		var i uint32
		var mask uint32 = timeNear << timeLevelShift
		for i = 0; i < 3; i++ {
			if (expire | (mask - 1)) == (currentTime | (mask - 1)) {
				break
			}
			mask <<= timeLevelShift
		}

		tm.t[i][(expire>>(timeNearShift+i*timeLevelShift))&timeLevelMask].PushBack(n)
	}
}

func (tm *Timer) putTime(t int, f TimeGoFunction, parm interface{}) {

	n := &Node{0, parm, f}
	tm.m.Lock()
	defer tm.m.Unlock()
	n.expire = uint32(t) + tm.timeDur
	tm.addNode(n)
}

func (tm *Timer) getTime() uint64 {
	return uint64(time.Now().UnixNano() / 10000000)
}
