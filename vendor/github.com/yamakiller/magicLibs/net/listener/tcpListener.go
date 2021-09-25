package listener

import (
	"net"
	"sync"
)

//SpawnTCPListener create an tcp listener
func SpawnTCPListener(l net.Listener) Listener {
	return &TCPListener{_l: l}
}

//TCPListener TCP listener
type TCPListener struct {
	_l  net.Listener
	_wg sync.WaitGroup
}

//Accept tcp accept connection
func (slf *TCPListener) Accept([]interface{}) (interface{}, error) {
	c, err := slf._l.Accept()
	if err != nil {
		return nil, err
	}

	slf._wg.Add(1)
	return &TCPConn{Conn: c, _wg: &slf._wg}, nil
}

//Addr Returns  address
func (slf *TCPListener) Addr() net.Addr {
	return slf._l.Addr()
}

//Close close listener
func (slf *TCPListener) Close() error {
	if err := slf._l.Close(); err != nil {
		return err
	}

	slf._wg.Wait()
	return nil
}

//ToString ....
func (slf *TCPListener) ToString() string {
	return "tcp listener"
}

//TCPConn tcp connection
type TCPConn struct {
	net.Conn
	_wg *sync.WaitGroup
}

//Close TCP 连接关闭
func (slf *TCPConn) Close() error {
	if slf._wg != nil {
		slf._wg.Done()
	}
	slf._wg = nil
	return slf.Conn.Close()
}
