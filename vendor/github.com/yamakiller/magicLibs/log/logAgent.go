package log

//New 创建一个日志对象
func New(f func() LogAgent) LogAgent {
	return f()
}

//LogAgent 日志对象接口
type LogAgent interface {
	Close()
	WithHandle(interface{})

	Out(fmrt string, args ...interface{})
	Info(prefix, fmrt string, args ...interface{})
	Error(prefix, fmrt string, args ...interface{})
	Debug(prefix, fmrt string, args ...interface{})
	Warning(prefix, fmrt string, args ...interface{})
	Trace(prefix, fmrt string, args ...interface{})
	Fatal(prefix, fmrt string, args ...interface{})
	Panic(prefix, fmrt string, args ...interface{})
}
