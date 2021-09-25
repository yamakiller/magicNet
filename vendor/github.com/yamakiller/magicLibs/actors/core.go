package actors

import (
	"fmt"
	"math"
	"sync"

	"github.com/yamakiller/magicLibs/actors/messages"

	"github.com/yamakiller/magicLibs/log"
)

//New 创建Actor内核
func New(sch Scheduler) *Core {
	c := &Core{
		_pidSets: &PIDSet{_pids: make([]PID, math.MaxUint16),
			_seq: 1,
		},
		_hs:  make(map[uint32]handle),
		_sch: sch,
	}

	c._deadLetter = &deathLetter{
		_parent: c,
		_sub:    make(chan interface{}, 32),
		_closed: make(chan bool),
	}

	c._deadLetter._wait.Add(1)
	go c._deadLetter.run()

	return c
}

//Core actors 核心模块
type Core struct {
	_pidSets    *PIDSet
	_log        log.LogAgent
	_hs         map[uint32]handle
	_gw         sync.WaitGroup
	_sch        Scheduler
	_deadLetter *deathLetter
	_syn        sync.Mutex
}

//New 创建一个Actor
func (slf *Core) New(f func(*PID) Actor) (*PID, error) {
	pid, err := slf._pidSets.Next()
	if err != nil {
		return nil, err
	}

	pid._parent = slf

	actor := f(pid)
	ctx := spawnContext(slf, actor, pid)
	hl := spawnHandle(ctx, slf._sch)

	slf._syn.Lock()
	if _, ok := slf._hs[pid.ID]; ok {
		slf._syn.Unlock()
		errID := pid.ID
		slf._pidSets.Remove(pid)

		return nil, fmt.Errorf("pid[[.%08x]] repeat", errID)
	}
	slf._hs[pid.ID] = hl
	slf._syn.Unlock()

	slf._gw.Add(1)
	pid.postSysMessage(messages.StartedMessage)

	return pid, nil
}

//Delete 删除一个actor
func (slf *Core) Delete(pid *PID) error {
	slf._syn.Lock()
	defer slf._syn.Unlock()
	id := pid.ID

	slf._pidSets.Remove(pid)
	if _, ok := slf._hs[id]; ok {
		delete(slf._hs, id)
	}

	slf._gw.Done()
	return nil
}

//WithLogger 设置日志接口
func (slf *Core) WithLogger(log log.LogAgent) {
	slf._log = log
}

//Close 关闭
func (slf *Core) Close() {

	slf._gw.Wait()

	if slf._deadLetter != nil {
		slf._deadLetter.close()
		slf._deadLetter = nil
	}
}

//getHandle 获取一个Actor Handle
func (slf *Core) getHandle(pid *PID) handle {
	slf._syn.Lock()
	defer slf._syn.Unlock()

	if h, ok := slf._hs[pid.ID]; ok {
		return h
	}
	return slf._deadLetter
}
