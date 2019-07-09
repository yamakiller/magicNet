package logger

type event struct {
	level   uint32
	prefix  string
	message string
}
