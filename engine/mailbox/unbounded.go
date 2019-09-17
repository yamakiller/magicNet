package mailbox

import (
	"github.com/yamakiller/magicNet/engine/overload"
	"github.com/yamakiller/magicNet/util"
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

// Unbounded  : 没有上限的邮箱制造器
func Unbounded(mailboxStats ...Statistics) Make {
	return func() Mailbox {
		q := &unboundedMailboxQueue{
			usrMailbox: overload.NewQueue(16),
		}

		return &defaultMailbox{
			sysMailbox:   util.NewQueue(),
			usrMailbox:   q,
			mailboxStats: mailboxStats,
		}
	}
}
