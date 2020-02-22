package ado

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/yamakiller/magicNet/netboxs"
)

type Conn struct {
	netboxs.BKCPConn
	_buffer []byte
}

func (slf *Conn) Keepalive() time.Duration {
	//return 4 * time.Second
	return 4 * time.Second
}

func (slf *Conn) Ping() {
	ping := "examples ping"
	slf.Push(ping)
}

func (slf *Conn) UnSeria() (interface{}, error) {
	n, err := slf.Reader().Recv(slf._buffer, int32(len(slf._buffer)))
	if err != nil {
		return nil, err
	}

	if n < 2 {
		return nil, errors.New("recv data size error min 2")
	}

	msg := ""
	header := binary.BigEndian.Uint16(slf._buffer)
	if int32(header) > (n - 2) {
		return nil, errors.New("recv data overflow")
	}

	if header > 0 {
		msg = string(slf._buffer[2 : header+2])
	}

	return msg, nil
}

func (slf *Conn) Seria(msg interface{}) error {
	length := len([]rune(msg.(string)))
	buffer := make([]byte, 2+length)
	binary.BigEndian.PutUint16(buffer, uint16(length))
	copy(buffer[2:], msg.(string))

	_, err := slf.Writer().Write(buffer, int32(2+length))
	if err != nil {
		return err
	}

	return nil
}

type ConnPools struct {
}

func (slf *ConnPools) Get() netboxs.Connect {
	return &Conn{BKCPConn: netboxs.BKCPConn{
		WriteQueueSize: 16,
	},
		_buffer: make([]byte, 2048)}
}

func (slf *ConnPools) Put(netboxs.Connect) {

}
