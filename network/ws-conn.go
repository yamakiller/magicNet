package network

import (
	"errors"

	"github.com/gorilla/websocket"
)

const (
	wsOutChanMax = 1024
)

type wsConn struct {
	sConn
}

func (sc *wsConn) getProto() string {
	return ProtoWeb
}

func (sc *wsConn) getType() int {
	return CConnect
}

func wsConnClose(s interface{}) {
	if conn, ok := s.(*websocket.Conn); ok {
		conn.Close()
	}
}

func wsConnRecv(s interface{}) (int, []byte, error) {
	conn, ok := s.(*websocket.Conn)
	if !ok {
		return 0, nil, errors.New("web conn object exception")
	}

	msgType, data, err := conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	if msgType != websocket.BinaryMessage {
		return 0, nil, errors.New("web conn receive only binary data")
	}

	return len(data), data, nil
}

func wsConnWrite(s interface{}, data []byte) (int, error) {
	conn, ok := s.(*websocket.Conn)
	if !ok {
		return 0, errors.New("web conn object exception")
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		return 0, err
	}

	return len(data), nil
}
