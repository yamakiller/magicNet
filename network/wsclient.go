package network

import (
	"magicNet/engine/actor"
	"sync"

	"github.com/gorilla/websocket"
)

type wsclient struct {
	h        int32
	c        *websocket.Conn
	w        sync.WaitGroup
	o        *actor.PID
	i        NetInfo
	out      chan *NetChunk
	outStat  int32  //out状态
	keepAive uint64 // 毫秒
	stat     int
}

func (wsc *wsclient) listen(operator *actor.PID, addr string) error {
	return nil
}

func (wsc *wsclient) connect(operator *actor.PID, addr string) error {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	wsc.c = c
	//启动读数据协程
	//启动写数据协程

	return nil
}

func (wsc *wsclient) read() {

}

func (wsc *wsclient) write() {

}
