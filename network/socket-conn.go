package network

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"
)

type sConn struct {
	_h        int32
	_s        interface{}
	_w        sync.WaitGroup
	_o        *actor.PID
	_i        NetInfo
	_rv       recvFunc
	_wr       writeFunc
	_cls      closeFunc
	_so       *slot
	_srv      *sServer
	_out      chan *NetChunk
	_quit     chan int
	_outStat  int32  //out状态
	_keepAive uint64 // 毫秒
	_stat     int32
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
	defer slf._w.Done()
	for {
		if slf._stat != Connecting &&
			slf._stat != Connected {
			goto read_end
		}

		n, data, err := slf._rv(slf._s)
		if err != nil {
			//? Log error log
			goto read_error
		}

		// Discard data
		if slf._stat != Connected {
			continue
		}

		slf._i.RecvBytes += uint64(n)
		slf._i.RecvLastTime = timer.Now()
		//Forwarding data  message
		actor.DefaultSchedulerContext.Send(slf._o, &NetChunk{Handle: slf._h, Data: data[:n]})
	}
read_error:
	slf._stat = Closing
	slf._cls(slf._s)
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)
	slf._so.l.Lock()
	closeHandle = slf._h
	closeOperator = slf._o
	close(slf._quit)

	//-----Waiting for the end of the write corout------
	ick := 0
	for {
		if atomic.CompareAndSwapInt32(&slf._outStat, 1, 1) {
			break
		}

		ick++
		if ick > 8 {
			ick = 0
			time.Sleep(time.Duration(10) * time.Millisecond)
		}
	}

	close(slf._out)
	//----------------------
	if slf._srv != nil {
		slf._srv._conns.Delete(slf._h)
	}

	slf._so.s = nil
	slf._so.b = resIdle
	slf._so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{Handle: closeHandle})
}

func (slf *sConn) write() {
	defer slf._w.Done()
	for {
		if slf._stat != Connecting &&
			slf._stat != Connected {
			goto write_end
		}

		select {
		case msg := <-slf._out:
			if slf._stat != Connecting &&
				slf._stat != Connected {
				goto write_end
			}

			n, err := slf._wr(slf._s, msg.Data)
			if err != nil {
				goto write_error
			}

			slf._i.WriteBytes += uint64(n)
			slf._i.WriteLastTime = timer.Now()
		case <-slf._quit:
			goto write_end
		}
	}
write_error:
	slf._stat = Closing
write_end:
	slf._outStat = 1
}

func (slf *sConn) push(data *NetChunk, n int) error {
	select {
	case <-slf._quit:
		return errors.New("conn closed")
	default:
	}

	select {
	case <-slf._quit:
		return errors.New("conn closed")
	case slf._out <- data:
	}

	return nil
}

func (slf *sConn) setKeepAive(keep uint64) {
	slf._keepAive = keep
}

func (slf *sConn) getKeepAive() uint64 {

	return slf._keepAive
}

func (slf *sConn) getLastActivedTime() uint64 {
	return slf._i.RecvLastTime
}

func (slf *sConn) getStat() int32 {
	return slf._stat
}

func (slf *sConn) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf._stat, Connecting, Connected)
}

func (slf *sConn) close(lck *mutex.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if slf._stat != Closing {
		slf._stat = Closing
		slf._cls(slf._s)
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (slf *sConn) closewait() {
	slf._w.Wait()
}
