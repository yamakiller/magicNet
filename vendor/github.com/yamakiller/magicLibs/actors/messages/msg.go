package messages

var (
	//StartedMessage ...
	StartedMessage interface{} = &Started{}
	//StopMessage ...
	StopMessage interface{} = &Stop{}
	//StoppingMessage ...
	StoppingMessage interface{} = &Stopping{}
	//StoppedMessage ...
	StoppedMessage interface{} = &Stopped{}
	//ResumeMessage ...
	ResumeMessage interface{} = Resume{}
	//SuspendMessage ...
	SuspendMessage interface{} = Suspend{}
)
