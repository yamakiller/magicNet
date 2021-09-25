package borker

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"

	"github.com/yamakiller/magicLibs/net/listener"
)

//UDPBorker UDP代理服务
type UDPBorker struct {
	Mtu      int
	OutQueue int

	Spawn   func(*listener.UDPReport) error
	_listen *listener.UDPListener
	_wg     sync.WaitGroup
	_closed chan bool
}

//ListenAndServe Listen and startup server
func (slf *UDPBorker) ListenAndServe(addr string) error {
	if slf._closed == nil {
		slf._closed = make(chan bool, 1)
	}

	address, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	lst, err := net.ListenUDP("udp", address)
	if err != nil {
		return err
	}

	slf._listen = listener.SpawnUDPListener(lst, slf.Mtu, slf.OutQueue)

	slf._wg.Add((2))
	go slf.Serve()
	go slf.update()
	return nil
}

func (slf *UDPBorker) ListenAndServeTls(addr string, ptls *tls.Config) error {
	return errors.New("undefined tls")
}

//Serve udp recv
func (slf *UDPBorker) Serve() error {
	defer func() {
		slf._wg.Done()
	}()

	var err error
	for {
		select {
		case <-slf._closed:
			goto exit
		default:
			c, e := slf._listen.Accept(nil)
			if e != nil {
				err = e
				goto exit
			}

			if c == nil {
				continue
			}

			if e := slf.Spawn(c.(*listener.UDPReport)); e != nil {
				c.(*listener.UDPReport).Close()
				continue
			}
		}
	}
exit:
	return err
}

func (slf *UDPBorker) update() {
	defer func() {
		slf._wg.Done()
	}()
	for {
		select {
		case <-slf._closed:
			goto exit
		default:
			_, err := slf._listen.Wait()
			if err != nil {
				goto exit
			}
		}
	}
exit:
}

//Listener Returns listenner object
func (slf *UDPBorker) Listener() listener.Listener {
	return slf._listen
}

//Shutdown 关闭服务
func (slf *UDPBorker) Shutdown() {
	slf._closed <- true
	if slf._listen != nil {
		slf._listen.Close()
	}
	slf._wg.Wait()

	slf._listen = nil
	close(slf._closed)
}
