package netboxs

import "time"

//Connect netboxs connectioner interface
type Connect interface {
	Socket() int32
	Keepalive() time.Duration
	WithSocket(int32)
	WithIO(interface{})
	//Write([]byte) error
	Ping()
	Parse() (interface{}, error)
	UnParse(interface{}) error
	Push(interface{}) error
	Pop() chan interface{}
	Close() error
}
