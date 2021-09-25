package borker

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/net/listener"
	"github.com/yamakiller/magicLibs/net/middle"
	"github.com/yamakiller/magicLibs/util"
)

//KCPBorker kcp 网络代理服务
type KCPBorker struct {
	RecvWndSize   int32
	SendWndSize   int32
	RecvQueueSize int32
	NoDelay       int32
	Interval      int32
	Resend        int32
	Nc            int32
	RxMinRto      int32
	FastResend    int32
	Mtu           int
	Middleware    middle.KSMiddleware

	Spawn func(*listener.KCPConn) error

	_listen *listener.KCPListener
	_wg     sync.WaitGroup
	_closed chan bool
}

//ListenAndServe 监听并启动服务
func (slf *KCPBorker) ListenAndServe(addr string) error {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	if slf.Mtu == 0 {
		slf.Mtu = 1400
	}

	slf.Mtu = int(mmath.Align(uint32(slf.Mtu), uint32(4)))

	slf._listen = listener.SpawnKCPListener(conn, slf.Mtu)
	slf._listen.RecvWndSize = slf.RecvWndSize
	slf._listen.SendWndSize = slf.SendWndSize
	slf._listen.RecvQueueSize = slf.RecvQueueSize
	slf._listen.NoDelay = slf.NoDelay
	slf._listen.Interval = slf.Interval
	slf._listen.Resend = slf.Resend
	slf._listen.Nc = slf.Nc
	slf._listen.RxMinRto = slf.RxMinRto
	slf._listen.FastResend = slf.FastResend
	slf._listen.Middleware = slf.Middleware

	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	slf._wg.Add((2))
	go slf.Serve()
	go slf.update()
	return nil
}

func (slf *KCPBorker) ListenAndServeTls(addr string, ptls *tls.Config) error {
	return errors.New("undefined tls")
}

//Serve 启动服务
func (slf *KCPBorker) Serve() error {
	defer func() {
		slf._wg.Done()
	}()

	var err error
	params := [1]interface{}{slf.Spawn}
	for {
		select {
		case <-slf._closed:
			goto exit
		default:
			_, e := slf._listen.Accept(params[:])
			if e != nil {
				err = e
				goto exit
			}

			/*if c == nil {
				continue
			}

			if e := slf.Spawn(c.(*listener.KCPConn)); e != nil {
				c.(*listener.KCPConn).Close()
				continue
			}*/
		}
	}
exit:
	return err
}

func (slf *KCPBorker) update() {
	defer func() {
		slf._wg.Done()
	}()
	for {
		select {
		case <-slf._closed:
			for {
				current := util.Timestamp()
				if n := slf._listen.Update(current); n < 0 {
					break
				}
				time.Sleep(time.Millisecond)
			}
			goto exit
		default:
			current := util.Timestamp()
			slf._listen.Update(current)
			time.Sleep(time.Duration(slf.Interval/2) * time.Millisecond)
		}
	}
exit:
}

//Listener Returns listenner object
func (slf *KCPBorker) Listener() listener.Listener {
	return slf._listen
}

//Shutdown 关闭服务
func (slf *KCPBorker) Shutdown() {
	if slf._closed != nil {
		select {
		case <-slf._closed:
		default:
			close(slf._closed)
		}
	}
	if slf._listen != nil {
		slf._listen.Close()
	}
	slf._wg.Wait()
	slf._listen = nil
}
