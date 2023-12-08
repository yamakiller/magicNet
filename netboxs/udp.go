package netboxs

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/net/borker"
	"github.com/yamakiller/magicLibs/net/listener"
	"github.com/yamakiller/magicNet/netmsgs"
)

// UDPBox udp network box
type UDPBox struct {
	Mtu        int
	OutQueue   int
	SeriaFun   func(interface{}) ([]byte, error)
	UnSeriaFun func([]byte) (interface{}, error)

	boxs.Box
	_borker *borker.UDPBorker
}

// WithPool Disable
func (slf *UDPBox) WithPool(pool Pool) {
	panic("UDP WithPool Function is disable")
}

// WithMax Disable
func (slf *UDPBox) WithMax(max int32) {
	panic("UDP WithMax Function is disable")
}

// ListenAndServe 启动监听服务
func (slf *UDPBox) ListenAndServe(addr string) error {
	slf.Box.StartedWait()
	slf._borker = &borker.UDPBorker{
		Spawn:    slf.handleConnect,
		Mtu:      slf.Mtu,
		OutQueue: slf.OutQueue,
	}

	if err := slf._borker.ListenAndServe(addr); err != nil {
		return err
	}

	return nil
}

func (slf *UDPBox) ListenAndServeTls(addr string, ptls *tls.Config) error {
	return errors.New("undefined listen tls")
}

func (slf *UDPBox) handleConnect(report *listener.UDPReport) error {
	msg, err := slf.UnSeriaFun(report.Recv())
	if err != nil {
		slf.Box.GetPID().Post(&netmsgs.Error{Sock: report.Addr(), Err: err})
		return nil
	}

	slf.Box.GetPID().Post(&netmsgs.Message{Sock: report.Addr(), Data: msg})
	return nil
}

// OpenTo Disable
func (slf *UDPBox) OpenTo(socket interface{}) error {
	panic("UDP OpenTo Function is disable")
}

// SendTo send data to address
// Param UDPAddr
// Param []byte/interface{}
func (slf *UDPBox) SendTo(addrs interface{}, msg interface{}) error {
	var (
		err error
		wby []byte
	)
	addr, ok := addrs.(*net.UDPAddr)
	if !ok {
		return errors.New("param (1):socket is UDPAddr")
	}

	if slf.SeriaFun == nil {
		return errors.New("need define seria function")
	}

	wby, err = slf.SeriaFun(msg)
	if err != nil {
		return err
	}

	if _, err = slf._borker.Listener().(*listener.UDPListener).WriteTo(*addr, wby); err != nil {
		return err
	}

	return nil
}

// CloseTo Disable
func (slf *UDPBox) CloseTo(socket int32) error {
	panic("UDP CloseTo Function is disable")
}

// CloseToWait Disable
func (slf *UDPBox) CloseToWait(socket int32) error {
	panic("UDP CloseToWait Function is disable")
}

// GetConnect Disable
func (slf *UDPBox) GetConnect(socket int32) (interface{}, error) {
	panic("UDP GetConnect Function is disable")
}

// GetValues Disable
func (slf *UDPBox) GetValues() []int32 {
	panic("UDP GetValues Function is disable")
}

// Shutdown 关闭服务
func (slf *UDPBox) Shutdown() {
	slf._borker.Shutdown()
	slf.Box.Shutdown()
	slf.Box.StoppedWait()
}
