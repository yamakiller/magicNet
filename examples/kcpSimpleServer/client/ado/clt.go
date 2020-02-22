package ado

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/yamakiller/magicLibs/net/connection"
	"github.com/yamakiller/mgokcp/mkcp"
)

func NewClt(id uint32) *Clt {
	c := &Clt{
		KCPClient: connection.KCPClient{
			WriteWaitQueue: 32,
			ReadWaitQueue:  32,
			S:              &CltSeria{},
			E:              &CltExption{},
		},
	}

	c.WithID(id)
	return c
}

type Clt struct {
	connection.KCPClient
	Timeout time.Time
	Check   int
}

func (slf *Clt) ReadLoop() {
	for {
		_, err := slf.Parse()
		if err != nil {
			slf.Close()
			break
		}
	}
}

type CltSeria struct {
}

func (slf *CltSeria) UnSeria(buf []byte) (interface{}, error) {
	if len(buf) < 2 {
		return nil, errors.New("data min 2")
	}

	n := len(buf)
	msg := ""
	header := binary.BigEndian.Uint16(buf)
	if int(header) > (n - 2) {
		return nil, errors.New("recv data overflow")
	}

	if header > 0 {
		msg = string(buf[2 : header+2])
	}

	return msg, nil
}

func (slf *CltSeria) Seria(msg interface{}, w *mkcp.KCP) (int, error) {
	length := len([]rune(msg.(string)))
	buffer := make([]byte, 2+length)
	binary.BigEndian.PutUint16(buffer, uint16(length))
	copy(buffer[2:], msg.(string))

	n, err := w.Send(buffer, int32(2+length))
	if err != nil {
		return 0, err
	}

	return int(n), nil
}

type CltExption struct {
}

func (slf *CltExption) Error(e error) {
	fmt.Println("Clt Exption:", e)
}
