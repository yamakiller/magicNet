package ado

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/yamakiller/magicLibs/net/connection"
)

type Clt struct {
	connection.WSSClient
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

func (slf *CltSeria) UnSeria(reader io.Reader) (interface{}, error) {
	var nlen uint16
	err := binary.Read(reader, binary.BigEndian, &nlen)
	if err != nil {
		return nil, err
	}

	buff := make([]byte, nlen)
	_, err = reader.Read(buff)
	if err != nil {
		return nil, err
	}
	return string(buff), nil
}

func (slf *CltSeria) Seria(msg interface{}, w io.Writer) (int, error) {
	data := msg.(string)
	length := uint16(len([]rune(data)))

	err := binary.Write(w, binary.BigEndian, &length)
	if err != nil {
		return -1, err
	}

	_, err = w.Write([]byte(data))
	if err != nil {
		return -1, err
	}

	return int(length + 2), nil
}
