package netboxs

import "time"

//Connect netboxs connectioner interface
type Connect interface {
	Socket() int32
	Keepalive() time.Duration
	WithSocket(int32)
	WithIO(interface{})
	Write([]byte) error
	Ping()
	Parse() (interface{}, error)
	Push([]byte) error
	Pop() chan []byte
	Close() error
}
