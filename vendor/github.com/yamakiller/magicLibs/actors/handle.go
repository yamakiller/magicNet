package actors

import (
	"sync/atomic"

	"github.com/yamakiller/magicLibs/actors/messages"
)

func spawnHandle(in invoker, sch Scheduler) handle {
	return &actorHandle{
		_mailbox: mailbox{
			_usrMailbox: spawnQueue(8),
			_sysMailbox: spawnQueue(4),
			_dispatcher: &dispatcher{
				_sch: sch,
			},
			_invoker: in,
		},
	}
}

type handle interface {
	overloadUsrMessage() int
	postUsrMessage(pid *PID, message interface{})
	postSysMessage(pid *PID, message interface{})
	Stop(pid *PID)
}

type actorHandle struct {
	_mailbox mailbox
	_death   int32
}

//overloadUsrMessage Returns user message overload warring
func (slf *actorHandle) overloadUsrMessage() int {
	return slf._mailbox._usrMailbox.Overload()
}

// postUsrMessage Send user level messages
func (slf *actorHandle) postUsrMessage(pid *PID, message interface{}) {
	slf._mailbox.postUsrMessage(message)
}

// postSysMessage Send system level messages
func (slf *actorHandle) postSysMessage(pid *PID, message interface{}) {
	slf._mailbox.postSysMessage(message)
}

// Stop Send stop Actor message
func (slf *actorHandle) Stop(pid *PID) {
	atomic.StoreInt32(&slf._death, 1)
	slf.postSysMessage(pid, messages.StopMessage)
}
