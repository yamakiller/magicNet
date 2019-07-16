package network

import (
	"magicNet/engine/actor"
	"magicNet/timer"
	"sync"
	"sync/atomic"

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
	return nil
}

func (wsc *wsclient) read(so *slot) {
	for {
		if wsc.stat != Connecting && wsc.stat != Connected {
			goto read_end
		}

		msgType, data, err := wsc.c.ReadMessage()
		if err != nil {
			//记录错误日志
			goto read_error
		}

		// 不接收非二进制编码数据
		if msgType != websocket.BinaryMessage {
			goto read_error
		}

		// 丢弃数据包
		if wsc.stat != Connected {
			continue
		}

		wsc.i.ReadBytes += uint64(len(data))
		wsc.i.ReadLastTime = timer.Now()
		//数据包丢给 Actor
		actor.DefaultSchedulerContext.Send(wsc.o, &NetChunk{Data: data})
	}

read_error:
	wsc.stat = Closing
	wsc.c.Close()
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	so.l.Lock()
	closeHandle = wsc.h
	closeOperator = wsc.o
	close(wsc.out)
	//-----等待写协程结束------
	for {
		if atomic.CompareAndSwapInt32(&wsc.outStat, 1, 1) {
			break
		}
	}
	so.s = nil
	so.b = resIdle
	so.l.Unlock()

	wsc.w.Done()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{handle: closeHandle})
}

func (wsc *wsclient) write(so *slot) {
	for {
		if wsc.stat != Connecting && wsc.stat != Connected {
			goto write_end
		}

		select {
		case msg := <-wsc.out:
			if wsc.stat != Connecting && wsc.stat != Connected {
				goto write_end
			}

			if err := wsc.c.WriteMessage(websocket.BinaryMessage, msg.Data); err != nil {
				goto write_error
			}

			wsc.i.WriteBytes += uint64(len(msg.Data))
			wsc.i.WriteLastTime = timer.Now()
		}
	}

write_error:
	wsc.stat = Closing
write_end:
	wsc.w.Done()
	wsc.outStat = 1
}
