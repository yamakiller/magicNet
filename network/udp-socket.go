package network

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"
)

type udpSocket struct {
	_h        int32
	_s        *net.UDPConn
	_i        NetInfo
	_so       *slot
	_operator *actor.PID
	_netWait  sync.WaitGroup
	_out      chan *NetChunk
	_quit     chan int
	_outStat  int32
	_mode     int32
	_stat     int32
}

func (slf *udpSocket) listen(operator *actor.PID, addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	slf._s = ln
	slf._mode = 0
	slf._stat = Connecting
	slf._netWait.Add(2)
	go slf.recv()
	go slf.write()

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (slf *udpSocket) connect(operator *actor.PID, addr string) error {
	return nil
}

func (slf *udpSocket) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	udpSrcAddr, _ := net.ResolveUDPAddr("udp", srcAddr)
	udpDstAddr, _ := net.ResolveUDPAddr("udp", dstAddr)

	ln, err := net.DialUDP("udp", udpSrcAddr, udpDstAddr)
	if err != nil {
		return err
	}
	slf._s = ln
	slf._mode = 1
	slf._stat = Connecting

	return nil
}

func (slf *udpSocket) recv() {
	defer slf._netWait.Done()
	for {
		if slf._stat != Connecting &&
			slf._stat != Connected {
			goto read_end
		}

		var inBuf []byte
		n, addr, err := slf._s.ReadFrom(inBuf)
		if err != nil {
			goto read_error
		}

		if slf._stat != Connected {
			continue
		}

		slf._i.RecvBytes += uint64(n)
		slf._i.RecvLastTime = timer.Now()

		udpAddr, _ := net.ResolveUDPAddr(addr.Network(), addr.String())

		actor.DefaultSchedulerContext.Send(slf._operator, &NetChunk{Handle: slf._h, Data: inBuf, Addr: udpAddr.IP, Port: uint16(udpAddr.Port)})
	}
read_error:
	slf._stat = Closing
	slf._s.Close()
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	slf._so.l.Lock()
	closeHandle = slf._h
	closeOperator = slf._operator
	close(slf._quit)
	//-----Waiting for the write coroutine to end------
	for {
		if atomic.CompareAndSwapInt32(&slf._outStat, 1, 1) {
			break
		}
	}
	close(slf._out)

	slf._so.s = nil
	slf._so.b = resIdle
	slf._so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{Handle: closeHandle})
}

func (slf *udpSocket) write() {
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

			udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprint(msg.Addr.String(), ":", msg.Port))
			n, err := slf._s.WriteToUDP(msg.Data, udpAddr)
			if err != nil {
				//TODO:  ?
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
	slf._netWait.Done()
	slf._outStat = 1
}

func (slf *udpSocket) push(data *NetChunk, n int) error {
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

func (slf *udpSocket) setKeepAive(keep uint64) {

}

func (slf *udpSocket) getKeepAive() uint64 {
	return 0
}

func (slf *udpSocket) getLastActivedTime() uint64 {
	return slf._i.RecvLastTime
}

func (slf *udpSocket) getStat() int32 {
	return slf._stat
}

func (slf *udpSocket) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf._stat, Connecting, Connected)
}

func (slf *udpSocket) close(lck *mutex.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if slf._stat != Closing {
		slf._stat = Closing
		slf._s.Close()
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (slf *udpSocket) closewait() {
	slf._netWait.Wait()
}

func (slf *udpSocket) getProto() string {
	return protoUDP
}

func (slf *udpSocket) getType() int {
	if slf._mode == 0 {
		return CListen
	}
	return CClient
}
