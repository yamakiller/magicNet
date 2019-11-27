package network

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	sConn
}

func (slf *wsClient) connect(operator *actor.PID, addr string) error {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	now := timer.Now()
	slf._s = c
	slf._stat = Connecting
	slf._rv = wsConnRecv
	slf._wr = wsConnWrite
	slf._cls = wsConnClose
	slf._i.RecvLastTime = now
	slf._i.WriteLastTime = now
	return nil
}

func (slf *wsClient) getProto() string {
	return protoWeb
}

func (slf *wsClient) getType() int {
	return CClient
}
