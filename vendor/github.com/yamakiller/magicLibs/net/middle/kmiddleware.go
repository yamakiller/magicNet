package middle

import (
	"net"
	"time"
)

//KSMiddleware KCP 服务端中间件
type KSMiddleware interface {
	Exception
	Subscribe([]byte, *net.UDPConn, *net.UDPAddr) (interface{}, error)
	UnSubscribe(uint32)
	Update()
}

//KCMiddleware KCP 客户端中间件
type KCMiddleware interface {
	Subscribe(*net.UDPConn, *net.UDPAddr, time.Duration) (interface{}, error)
}
