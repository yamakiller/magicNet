package logger

// Event : 日志事件
type Event struct {
	level   uint32
	prefix  string
	message string
}
