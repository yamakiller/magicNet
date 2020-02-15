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
	_queue chan []byte
}

func (slf *Conn) Keepalive() time.Duration {
	//return 4 * time.Second
	return 0
}

func (slf *Conn) WithIO(c interface{}) {
	slf.BTCPConn.WithIO(c)
	slf._queue = make(chan []byte, 32)
}

func (slf *Conn) Ping() {
	ping := "examples ping"
	pingLen := len([]rune(ping))
	bs := make([]byte, 2)
	binary.BigEndian.PutUint16(bs, uint16(pingLen))
	bs = append(bs, ([]byte(ping))...)
	slf.Write(bs)
}

func (slf *Conn) Parse() (interface{}, error) {
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

func (slf *Conn) Push(data []byte) error {
	slf._queue <- data
	return nil
}

func (slf *Conn) Pop() chan []byte {
	return slf._queue
}

func (slf *Conn) Close() error {
	close(slf._queue)
	return nil
}

type ConnPools struct {
}

func (slf *ConnPools) Get() netboxs.Connect {
	return &Conn{BTCPConn: netboxs.BTCPConn{
		ReadBufferSize:  8192,
		WriteBufferSize: 8192,
	}}
}

func (slf *ConnPools) Put(netboxs.Connect) {

}
