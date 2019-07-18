package network

import (
	"errors"
	"net"
	"time"
)

type tcpConn struct {
	sConn
}

func (tpc *tcpConn) setKeepAive(keep uint64) {
	tpc.keepAive = keep
	if conn, ok := tpc.s.(*net.TCPConn); ok {
		if tpc.keepAive > 0 {
			conn.SetKeepAlive(true)
			conn.SetKeepAlivePeriod(time.Duration(tpc.keepAive) * time.Millisecond)
		} else {
			conn.SetKeepAlive(false)
			conn.SetKeepAlivePeriod(time.Duration(tpc.keepAive))
		}
	}
}

func (tpc *tcpConn) getProto() string {
	return protoTCP
}

func (tpc *tcpConn) getType() int {
	return CConnect
}

func tcpConnClose(s interface{}) {
	if conn, ok := s.(*net.TCPConn); ok {
		conn.Close()
	}
}

func tcpConnRecv(s interface{}) (int, []byte, error) {
	conn, ok := s.(*net.TCPConn)
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
	conn, ok := s.(*net.TCPConn)
	if !ok {
		return 0, errors.New("socket conn object exception")
	}

	n, err := conn.Write(data)
	if err != nil {
		return 0, err
	}

	return n, nil
}
