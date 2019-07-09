package mailbox

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月04日 20:34:38
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月09日 10:08:39
 * @Description: Actor 的消息队列模块
 */

import (
	"magicNet/engine/util"
	"runtime"
	"sync/atomic"
)

// Statistics : 统计接口
type Statistics interface {
	MailboxStarted()
	MessagePosted(message interface{})
	MessageReceived(message interface{})
	MailboxEmpty()
}

// MessageInvoker : 消息调度器接口
type MessageInvoker interface {
	InvokeUsrMessage(interface{})
	InvokeSysMessage(interface{})
	EscalateFailure(reason interface{}, message interface{})
}

// Mailbox : 消息邮箱列接口[消息队]
type Mailbox interface {
	PostUsrMessage(message interface{})
	PostSysMessage(message interface{})
	OverloadUsrMessage() int
	RegisterHandlers(invoker MessageInvoker, dispatcher Dispatcher)
	Start()
}

// Make : 邮箱制造器接口[消息队列]
type Make func() Mailbox

const (
	idle int32 = iota
	running
)

type defaultMailbox struct {
	usrMailbox      queue
	sysMailbox      *util.Queue
	schedulerStatus int32
	usrMessages     int32
	sysMessages     int32
	invoker         MessageInvoker
	dispatcher      Dispatcher
	suspended       int32
	mailboxStats    []Statistics
}

func (m *defaultMailbox) PostUsrMessage(message interface{}) {
	for _, ms := range m.mailboxStats {
		ms.MessagePosted(message)
	}

	m.usrMailbox.Push(message)
	atomic.AddInt32(&m.usrMessages, 1)
	m.schedule()
}

func (m *defaultMailbox) PostSysMessage(message interface{}) {
	for _, ms := range m.mailboxStats {
		ms.MessagePosted(message)
	}
	m.sysMailbox.Push(message)
	atomic.AddInt32(&m.sysMessages, 1)
	m.schedule()
}

func (m *defaultMailbox) OverloadUsrMessage() int {
	return m.usrMailbox.Overload()
}

func (m *defaultMailbox) RegisterHandlers(invoker MessageInvoker, dispatcher Dispatcher) {
	m.invoker = invoker
	m.dispatcher = dispatcher
}

func (m *defaultMailbox) schedule() {
	if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
		m.dispatcher.Schedule(m.processMessages)
	}
}

func (m *defaultMailbox) processMessages() {
process_lable:
	m.run()

	atomic.StoreInt32(&m.schedulerStatus, idle)
	sys := atomic.LoadInt32(&m.sysMessages)
	usr := atomic.LoadInt32(&m.usrMessages)

	if sys > 0 || (atomic.LoadInt32(&m.suspended) == 0 && usr > 0) {
		if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
			goto process_lable
		}
	}

	for _, ms := range m.mailboxStats {
		ms.MailboxStarted()
	}
}

func (m *defaultMailbox) run() {
	var msg interface{}

	//异常处理--------------- begin
	defer func() {
		if r := recover(); r != nil {
			m.invoker.EscalateFailure(r, msg)
		}
	}()
	//----------------------- end

	i, t := 0, m.dispatcher.Throughput()
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}

		i++
		if msg = m.sysMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.sysMessages, -1)
			switch msg.(type) {
			case *SuspendMailbox:
				atomic.StoreInt32(&m.suspended, 1)
			case *ResumeMailbox:
				atomic.StoreInt32(&m.suspended, 0)
			default:
				m.invoker.InvokeSysMessage(msg)
			}
			for _, ms := range m.mailboxStats {
				ms.MessageReceived(msg)
			}
			continue
		}

		if atomic.LoadInt32(&m.suspended) == 1 {
			return
		}

		if msg = m.usrMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.usrMessages, -1)
			m.invoker.InvokeUsrMessage(msg)
			for _, ms := range m.mailboxStats {
				ms.MessageReceived(msg)
			}
		} else {
			return
		}

	}

}

func (m *defaultMailbox) Start() {
	for _, ms := range m.mailboxStats {
		ms.MailboxStarted()
	}
}
