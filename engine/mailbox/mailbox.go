package mailbox

import (
  "runtime"
	"sync/atomic"
  "magicNet/engine/util"
)

type Statistics interface {
  MailboxStarted()
  MessagePosted(message interface{})
  MessageReceived(message interface{})
  MailboxEmpty()
}

type MessageInvoker interface {
  InvokeUsrMessage(interface{})
  InvokeSysMessage(interface{})
  EscalateFailure(reason interface{}, message interface{})
}

type Mailbox interface{
  PostUsrMessage(message interface{})
  PostSysMessage(message interface{})
  OverloadUsrMessage() int
  RegisterHandlers(invoker MessageInvoker, dispatcher Dispatcher)
  Start()
}

type Producer func() Mailbox

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
      for _, ms :=  range m.mailboxStats {
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

/*type Mailbox struct {
  usrBox chan interface{}
}

func NewMailbox()*Mailbox {
  return &Mailbox{make(chan interface{})}
}

func (m *Mailbox) PostUserMessage(message interface{}) {
  m.usrBox <- message
}*/