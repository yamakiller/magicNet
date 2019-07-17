package network

import (
	"errors"
	"net"
)

type tcpConn struct {
	sConn
}

func tcpConnClose(s interface{}) {
	if conn, ok := s.(net.Conn); ok {
		conn.Close()
	}
}

func tcpConnRecv(s interface{}) (int, []byte, error) {
	conn, ok := s.(net.Conn)
	if !ok {
		return 0, nil, errors.New("socket conn object exception")
	}
	var inBuf []byte
	n, err := conn.Read(inBuf)
	if err != nil {
		return 0, nil, err
	}

	return n, inBuf, nil
}

func tcpConnWrite(s interface{}, data []byte) (int, error) {
	conn, ok := s.(net.Conn)
	if !ok {
		return 0, errors.New("socket conn object exception")
	}

	n, err := conn.Write(data)
	if err != nil {
		return 0, err
	}

	return n, nil
}
