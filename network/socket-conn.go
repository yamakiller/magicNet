package network

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"
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
	quit     chan int
	outStat  int32  //out状态
	keepAive uint64 // 毫秒
	stat     int32
}

type recvFunc func(interface{}) (int, []byte, error)
type writeFunc func(interface{}, []byte) (int, error)
type closeFunc func(interface{})

func (slf *sConn) listen(operator *actor.PID, addr string) error {
	return nil
}

func (slf *sConn) connect(operator *actor.PID, addr string) error {
	return nil
}

func (slf *sConn) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (slf *sConn) recv() {
	defer slf.w.Done()
	for {
		if slf.stat != Connecting && slf.stat != Connected {
			goto read_end
		}

		n, data, err := slf.rv(slf.s)
		if err != nil {
			//? Log error log
			goto read_error
		}

		// Discard data
		if slf.stat != Connected {
			continue
		}

		slf.i.ReadBytes += uint64(n)
		slf.i.ReadLastTime = timer.Now()
		//Forwarding data  message
		actor.DefaultSchedulerContext.Send(slf.o, &NetChunk{Handle: slf.h, Data: data[:n]})
	}
read_error:
	slf.stat = Closing
	slf.cls(slf.s)
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)
	slf.so.l.Lock()
	closeHandle = slf.h
	closeOperator = slf.o
	close(slf.quit)

	//-----Waiting for the end of the write corout------
	for {
		if atomic.CompareAndSwapInt32(&slf.outStat, 1, 1) {
			break
		}
	}

	close(slf.out)
	//----------------------
	if slf.srv != nil {
		slf.srv.conns.Delete(slf.h)
	}

	slf.so.s = nil
	slf.so.b = resIdle
	slf.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{Handle: closeHandle})
}

func (slf *sConn) write() {
	defer slf.w.Done()
	for {
		if slf.stat != Connecting && slf.stat != Connected {
			goto write_end
		}

		select {
		case msg := <-slf.out:
			if slf.stat != Connecting && slf.stat != Connected {
				goto write_end
			}

			n, err := slf.wr(slf.s, msg.Data)
			if err != nil {
				goto write_error
			}

			slf.i.WriteBytes += uint64(n)
			slf.i.WriteLastTime = timer.Now()
		case <-slf.quit:
			goto write_end
		}
	}
write_error:
	slf.stat = Closing
write_end:
	slf.outStat = 1
}

func (slf *sConn) push(data *NetChunk, n int) error {
	select {
	case <-slf.quit:
		return errors.New("conn closed")
	default:
	}

	select {
	case <-slf.quit:
		return errors.New("conn closed")
	case slf.out <- data:
	}

	return nil
}

func (slf *sConn) setKeepAive(keep uint64) {
	slf.keepAive = keep
}

func (slf *sConn) getKeepAive() uint64 {

	return slf.keepAive
}

func (slf *sConn) getLastActivedTime() uint64 {
	return slf.i.ReadLastTime
}

func (slf *sConn) getStat() int32 {
	return slf.stat
}

func (slf *sConn) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf.stat, Connecting, Connected)
}

func (slf *sConn) close(lck *mutex.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if slf.stat != Closing {
		slf.stat = Closing
		slf.cls(slf.s)
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (slf *sConn) closewait() {
	slf.w.Wait()
}
