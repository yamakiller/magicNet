package borker

import (
	"crypto/tls"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/net/listener"
)

//TCPBorker tcp 网络代理服务
type TCPBorker struct {
	Spawn   func(net.Conn) error
	_listen listener.Listener
	_wg     sync.WaitGroup
	_closed chan bool
}

//ListenAndServe 监听并启动服务
func (slf *TCPBorker) ListenAndServe(addr string) error {
	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	lst, err := net.ListenTCP("tcp", address)

	if err != nil {
		return err
	}

	slf._listen = listener.SpawnTCPListener(lst)

	slf._wg.Add((1))
	go slf.Serve()
	return nil
}

func (slf *TCPBorker) ListenAndServeTls(addr string, ptls *tls.Config) error {
	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	lst, err := tls.Listen("tcp", addr, ptls)

	if err != nil {
		return err
	}

	slf._listen = listener.SpawnTCPListener(lst)

	slf._wg.Add((1))
	go slf.Serve()
	return nil
}

//Serve accept services
func (slf *TCPBorker) Serve() error {
	defer func() {
		slf._wg.Done()
	}()

	var err error
	tmpDelay := time.Duration(1) * time.Millisecond
	for {
		select {
		case <-slf._closed:
			goto exit
		default:
			c, e := slf._listen.Accept(nil)
			if e != nil {
				if v, ok := c.(*net.TCPConn); ok {
					v.SetNoDelay(true)
					v.SetKeepAlive(true)
				}

				if ne, ok := e.(net.Error); ok && ne.Temporary() {
					tmpDelay *= 5
					if max := 1 * time.Second; tmpDelay > max {
						tmpDelay = max
					}

					time.Sleep(tmpDelay)
					continue
				}

				if strings.Contains(e.Error(), "use of closed network connection") {
					continue
				}

				err = e
				goto exit
			}

			tmpDelay = time.Duration(1) * time.Millisecond
			if e := slf.Spawn(c.(net.Conn)); e != nil {
				c.(*listener.TCPConn).Close()
				continue
			}
		}
	}
exit:
	return err
}

//Listener Returns listenner object
func (slf *TCPBorker) Listener() listener.Listener {
	return slf._listen
}

//Shutdown 关闭服务
func (slf *TCPBorker) Shutdown() {
	slf._closed <- true
	if slf._listen != nil {
		slf._listen.Close()
	}
	slf._wg.Wait()
	slf._listen = nil
	close(slf._closed)
}
