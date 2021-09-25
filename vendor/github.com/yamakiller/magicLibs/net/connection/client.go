package connection

import (
	"io"
	"time"
)

//Client 客户端接口
type Client interface {
	Connect(addr string, timeout time.Duration) error
	Parse() (interface{}, error)
	SendTo(interface{}) error
	Close() error
}

//Exception 异常处理接口
type Exception interface {
	Error(error)
}

//Serialization 序列化反序列化接口
type Serialization interface {
	UnSeria(io.Reader) (interface{}, int, error)
	Seria(interface{}, io.Writer) (int, error)
}
