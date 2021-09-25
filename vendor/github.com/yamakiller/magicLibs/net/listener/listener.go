package listener

import "net"

//Listener 监听接口
type Listener interface {
	Addr() net.Addr
	Accept([]interface{}) (interface{}, error)
	Close() error
	ToString() string
}
