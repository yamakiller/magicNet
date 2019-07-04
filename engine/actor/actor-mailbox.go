package actor

type Mailbox struct {
  box chan interface{}
}

func NewMailbox()*Mailbox {
  return &Mailbox{make(chan interface{})}
}

func (m *Mailbox) PostMessage(message interface{}) {
  m.box <- message
}

func (m *Mailbox) processMessages()

func (m *Mailbox) run() {

}
