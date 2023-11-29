package ado

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/yamakiller/magicNet/netboxs"
)

type Conn struct {
	netboxs.BTCPConn
}

func (slf *Conn) Keepalive() time.Duration {
	//return 4 * time.Second
	return 0
}

func (slf *Conn) Ping() {
	ping := "examples ping"
	slf.Push(ping)
}

// UnSeria　反序列化
func (slf *Conn) UnSeria() (interface{}, error) {
	var header uint16
	if err := binary.Read(slf.Reader(), binary.BigEndian, &header); err != nil {
		return nil, err
	}

	sz := int(header)
	buffer := make([]byte, sz)
	offset := 0

	for offset < sz {
		if offset > sz {
			return nil, errors.New("something went to wrong(offset overs remianing length)")
		}

		i, err := slf.Reader().Read(buffer[offset:])
		offset += i
		if err != nil && offset < sz {
			// if we read whole size of message, ignore error at this time.
			return nil, err
		}
	}

	if int(header) > len(buffer) {
		return nil, fmt.Errorf("publish length: %d, buffer: %d", header, len(buffer))
	}

	return string(buffer[0:header]), nil
}

// Seria 序列化
func (slf *Conn) Seria(msg interface{}) error {
	length := len([]rune(msg.(string)))
	buffer := make([]byte, 2+length)
	binary.BigEndian.PutUint16(buffer, uint16(length))
	copy(buffer[2:], msg.(string))

	return netboxs.WriteToBuffer(slf.Writer(), buffer)
}

type ConnPools struct {
}

func (slf *ConnPools) Get() netboxs.Connect {
	return &Conn{BTCPConn: netboxs.BTCPConn{
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
		WriteQueueSize:  16,
	}}
}

func (slf *ConnPools) Put(netboxs.Connect) {

}
