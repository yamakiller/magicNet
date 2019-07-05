package actor

type AutoReceiveMessage interface {
	AutoReceiveMessage()
}

type NotInfluenceReceiveTimeout interface {
	NotInfluenceReceiveTimeout()
}

type SystemMessage interface {
	SystemMessage()
}

type ReceiveTimeout struct {}

type Restarting struct {}

type Stopping struct{}

type Stopped struct{}

type Started struct {}

type Restart struct {}

func (*Restarting) AutoReceiveMessage() {}
func (*Stopping) AutoReceiveMessage()   {}
func (*Stopped) AutoReceiveMessage()    {}
//func (*PoisonPill) AutoReceiveMessage() {}

//func (*Started) SystemMessage()      {}
func (*Stop) SystemMessage()         {}
//func (*Watch) SystemMessage()        {}
//func (*Unwatch) SystemMessage()      {}
//func (*Terminated) SystemMessage()   {}
//func (*Failure) SystemMessage()      {}
//func (*Restart) SystemMessage()      {}
//func (*continuation) SystemMessage() {}*/

var (
	restartingMessage     interface{} = &Restarting{}
	stoppingMessage       interface{} = &Stopping{}
	stoppedMessage        interface{} = &Stopped{}
	//poisonPillMessage     interface{} = &PoisonPill{}
	receiveTimeoutMessage interface{} = &ReceiveTimeout{}
)

var (
	//restartMessage        interface{} = &Restart{}
	//startedMessage        interface{} = &Started{}
	stopMessage           interface{} = &Stop{}
	//resumeMailboxMessage  interface{} = &mailbox.ResumeMailbox{}
	//suspendMailboxMessage interface{} = &mailbox.SuspendMailbox{}
)
