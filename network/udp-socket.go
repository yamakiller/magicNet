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
	h        int32
	s        *net.UDPConn
	i        NetInfo
	so       *slot
	operator *actor.PID
	netWait  sync.WaitGroup
	out      chan *NetChunk
	quit     chan int
	outStat  int32
	mode     int32
	stat     int32
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

	slf.s = ln
	slf.mode = 0
	slf.stat = Connecting
	slf.netWait.Add(2)
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
	slf.s = ln
	slf.mode = 1
	slf.stat = Connecting

	return nil
}

func (slf *udpSocket) recv() {
	defer slf.netWait.Done()
	for {
		if slf.stat != Connecting && slf.stat != Connected {
			goto read_end
		}

		var inBuf []byte
		n, addr, err := slf.s.ReadFrom(inBuf)
		if err != nil {
			goto read_error
		}

		if slf.stat != Connected {
			continue
		}

		slf.i.ReadBytes += uint64(n)
		slf.i.ReadLastTime = timer.Now()

		udpAddr, _ := net.ResolveUDPAddr(addr.Network(), addr.String())

		actor.DefaultSchedulerContext.Send(slf.operator, &NetChunk{Handle: slf.h, Data: inBuf, Addr: udpAddr.IP, Port: uint16(udpAddr.Port)})
	}
read_error:
	slf.stat = Closing
	slf.s.Close()
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	slf.so.l.Lock()
	closeHandle = slf.h
	closeOperator = slf.operator
	close(slf.quit)
	//-----等待写协程结束------
	for {
		if atomic.CompareAndSwapInt32(&slf.outStat, 1, 1) {
			break
		}
	}
	close(slf.out)

	slf.so.s = nil
	slf.so.b = resIdle
	slf.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{Handle: closeHandle})
}

func (slf *udpSocket) write() {
	for {
		if slf.stat != Connecting && slf.stat != Connected {
			goto write_end
		}

		select {
		case msg := <-slf.out:
			if slf.stat != Connecting && slf.stat != Connected {
				goto write_end
			}

			udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprint(msg.Addr.String(), ":", msg.Port))
			n, err := slf.s.WriteToUDP(msg.Data, udpAddr)
			if err != nil {
				//?
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
	slf.netWait.Done()
	slf.outStat = 1
}

func (slf *udpSocket) push(data *NetChunk, n int) error {
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

func (slf *udpSocket) setKeepAive(keep uint64) {

}

func (slf *udpSocket) getKeepAive() uint64 {
	return 0
}

func (slf *udpSocket) getLastActivedTime() uint64 {
	return slf.i.ReadLastTime
}

func (slf *udpSocket) getStat() int32 {
	return slf.stat
}

func (slf *udpSocket) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf.stat, Connecting, Connected)
}

func (slf *udpSocket) close(lck *mutex.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if slf.stat != Closing {
		slf.stat = Closing
		slf.s.Close()
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (slf *udpSocket) closewait() {
	slf.netWait.Wait()
}

func (slf *udpSocket) getProto() string {
	return protoUDP
}

func (slf *udpSocket) getType() int {
	if slf.mode == 0 {
		return CListen
	}
	return CClient
}
