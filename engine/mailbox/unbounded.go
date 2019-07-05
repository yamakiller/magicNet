package mailbox

import (
  "magicNet/engine/overload"
  "magicNet/engine/util"
)

type unboundedMailboxQueue struct {
  usrMailbox *overload.Queue
}

func (q *unboundedMailboxQueue) Push(m interface{}) {
  q.usrMailbox.Push(m)
}

func (q *unboundedMailboxQueue) Pop() interface{} {
  m, o := q.usrMailbox.Pop()
  if o {
    return m
  }
  return nil
}

func (q *unboundedMailboxQueue) Overload() int {
  return q.usrMailbox.Overload()
}

func Unbounded(mailboxStats ...Statistics) Producer {
  return func() Mailbox {
    q := &unboundedMailboxQueue{
      usrMailbox: overload.NewQueue(16),
    }

    return &defaultMailbox {
      sysMailbox: util.NewQueue(),
      usrMailbox: q,
      mailboxStats: mailboxStats,
    }
  }
}
