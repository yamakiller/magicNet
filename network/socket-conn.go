package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/util"
	"magicNet/timer"
	"sync"
	"sync/atomic"
)

type sConn struct {
	h        int32
	s        interface{}
	w        sync.WaitGroup
	o        *actor.PID
	i        NetInfo
	rv       recvFunc
	wr       writeFunc
	cls      closeFunc
	so       *slot
	srv      *sServer
	out      chan *NetChunk
	outStat  int32  //out状态
	keepAive uint64 // 毫秒
	stat     int32
}

type recvFunc func(interface{}) (int, []byte, error)
type writeFunc func(interface{}, []byte) (int, error)
type closeFunc func(interface{})

func (sc *sConn) listen(operator *actor.PID, addr string) error {
	return nil
}

func (sc *sConn) connect(operator *actor.PID, addr string) error {
	return nil
}

func (sc *sConn) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (sc *sConn) recv() {
	defer sc.w.Done()
	for {
		if sc.stat != Connecting && sc.stat != Connected {
			goto read_end
		}

		n, data, err := sc.rv(sc.s)
		if err != nil {
			//记录错误日志
			goto read_error
		}

		// 丢弃数据包
		if sc.stat != Connected {
			continue
		}

		sc.i.ReadBytes += uint64(n)
		sc.i.ReadLastTime = timer.Now()
		//数据包丢给 Actor
		actor.DefaultSchedulerContext.Send(sc.o, &NetChunk{Data: data})
	}
read_error:
	sc.stat = Closing
	sc.cls(sc.s)
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	sc.so.l.Lock()
	closeHandle = sc.h
	closeOperator = sc.o
	close(sc.out)
	//-----等待写协程结束------
	for {
		if atomic.CompareAndSwapInt32(&sc.outStat, 1, 1) {
			break
		}
	}
	//----------------------
	if sc.srv != nil {
		sc.srv.conns.Delete(sc.h)
	}

	sc.so.s = nil
	sc.so.b = resIdle
	sc.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, NetClose{Handle: closeHandle})
}

func (sc *sConn) write() {
	for {
		if sc.stat != Connecting && sc.stat != Connected {
			goto write_end
		}

		select {
		case msg := <-sc.out:
			if sc.stat != Connecting && sc.stat != Connected {
				goto write_end
			}

			n, err := sc.wr(sc.s, msg.Data)
			if err != nil {
				goto write_error
			}

			sc.i.WriteBytes += uint64(n)
			sc.i.WriteLastTime = timer.Now()
		}
	}
write_error:
	sc.stat = Closing
write_end:
	sc.w.Done()
	sc.outStat = 1
}

func (sc *sConn) push(data *NetChunk, n int) error {
	//? 是否可以优化
	sc.out <- data
	return nil
}

func (sc *sConn) setKeepAive(keep uint64) {
	sc.keepAive = keep
}

func (sc *sConn) getKeepAive() uint64 {

	return sc.keepAive
}

func (sc *sConn) getLastActivedTime() uint64 {
	return sc.i.ReadLastTime
}

func (sc *sConn) getStat() int32 {
	return sc.stat
}

func (sc *sConn) setConnected() bool {
	return atomic.CompareAndSwapInt32(&sc.stat, Connecting, Connected)
}

func (sc *sConn) close(lck *util.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if sc.stat != Closing {
		sc.stat = Closing
		sc.cls(sc.s)
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (sc *sConn) closewait() {
	sc.w.Wait()
}
