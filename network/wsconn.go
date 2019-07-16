package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/util"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	wsOutChanMax = 1024
)

type wsconn struct {
	h        int32
	s        *websocket.Conn
	w        sync.WaitGroup
	o        *actor.PID
	i        NetInfo
	out      chan *NetChunk
	outStat  int32  //out状态
	keepAive uint64 // 毫秒
	stat     int
}

func (wsc *wsconn) listen(operator *actor.PID, addr string) error {
	return nil
}

func (wsc *wsconn) connect(operator *actor.PID, addr string) error {
	return nil
}

func (wsc *wsconn) setKeepAive(keep uint64) {
	wsc.keepAive = keep
}

func (wsc *wsconn) getStat() int {
	return wsc.stat
}

func (wsc *wsconn) setStat(stat int) {
	wsc.stat = stat
}

func (wsc *wsconn) close(lck *util.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if wsc.stat != Closing {
		wsc.stat = Closing
		wsc.s.Close()
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (wsc *wsconn) closewait() {
	wsc.w.Wait()
}
