package network

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"

	"github.com/gorilla/websocket"
)

type wsClient struct {
	sConn
}

func (wsc *wsClient) connect(operator *actor.PID, addr string) error {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	now := timer.Now()
	wsc.s = c
	wsc.stat = Connecting
	wsc.rv = wsConnRecv
	wsc.wr = wsConnWrite
	wsc.cls = wsConnClose
	wsc.i.ReadLastTime = now
	wsc.i.WriteLastTime = now
	return nil
}

func (wsc *wsClient) getProto() string {
	return protoWeb
}

func (wsc *wsClient) getType() int {
	return CClient
}
