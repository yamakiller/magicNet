package listener

import (
	"bytes"
	"errors"
	"net"
	"sync"
)

//SpawnUDPListener create an udp listener
func SpawnUDPListener(l *net.UDPConn, mtu, out int) *UDPListener {
	return &UDPListener{_l: l,
		_o:      make(chan *UDPTo, out),
		_closed: make(chan bool, 1),
		_mtu:    mtu}
}

//UDPListener UDP Listener
type UDPListener struct {
	_l      *net.UDPConn
	_o      chan *UDPTo
	_w      sync.WaitGroup
	_closed chan bool
	_mtu    int
}

//Accept udp message report
func (slf *UDPListener) Accept([]interface{}) (interface{}, error) {
	b := make([]byte, slf._mtu)
	n, addr, err := slf._l.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}

	return &UDPReport{_addr: *addr, _message: b, _length: n}, nil
}

//Wait ...
func (slf *UDPListener) Wait() (int, error) {

	slf._w.Add(1)
	defer slf._w.Done()

	select {
	case <-slf._closed:
		break
	case d, ok := <-slf._o:
		if !ok {
			return -1, nil
		}
		return slf._l.WriteToUDP(d._message, &d._addr)
	}
	return 0, errors.New("closed udp listener")
}

//Addr Returns  address
func (slf *UDPListener) Addr() net.Addr {
	return slf._l.LocalAddr()
}

//WriteTo Send to Udp Target address
func (slf *UDPListener) WriteTo(addr net.UDPAddr, buffer []byte) (int32, error) {
	slf._w.Add(1)
	defer slf._w.Done()

	size := len(buffer)

	select {
	case <-slf._closed:
		return 0, errors.New("udp listener closed")
	default:
		b := bytes.NewBuffer([]byte{})
		if n, err := b.Write(buffer); err != nil {
			return int32(n), err
		}
		slf._o <- &UDPTo{_addr: addr, _message: b.Bytes()}
	}

	return int32(size), nil
}

//Close close listener
func (slf *UDPListener) Close() error {
	slf._closed <- true
	if err := slf._l.Close(); err != nil {
		return err
	}
	slf._w.Wait()
	close(slf._closed)
	close(slf._o)

	return nil
}

//ToString ...
func (slf *UDPListener) ToString() string {
	return "udp listener"
}

//UDPTo udp send packet
type UDPTo struct {
	_addr    net.UDPAddr //source udp address
	_message []byte      //data message
}

//UDPReport tcp connection
type UDPReport struct {
	_addr    net.UDPAddr //source udp address
	_message []byte      //data message
	_length  int
}

//Addr udp address
func (slf *UDPReport) Addr() *net.UDPAddr {
	return &slf._addr
}

//Recv recvice data
func (slf *UDPReport) Recv() []byte {
	return slf._message[:slf._length]
}

//Close ...
func (slf *UDPReport) Close() error {
	return nil
}
