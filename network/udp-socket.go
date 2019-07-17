package network

import (
	"fmt"
	"magicNet/engine/actor"
	"magicNet/engine/util"
	"magicNet/timer"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type udpSocket struct {
	h        int32
	s        *net.UDPConn
	i        NetInfo
	so       *slot
	operator *actor.PID
	netWait  sync.WaitGroup
	out      chan *NetChunk
	outStat  int32
	mode     int32
	stat     int32
}

func (ups *udpSocket) listen(operator *actor.PID, addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}

	ups.s = ln
	ups.mode = 0
	ups.stat = Connecting
	ups.netWait.Add(2)
	go ups.recv()
	go ups.write()

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (ups *udpSocket) connect(operator *actor.PID, addr string) error {
	return nil
}

func (ups *udpSocket) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	udpSrcAddr, _ := net.ResolveUDPAddr("udp", srcAddr)
	udpDstAddr, _ := net.ResolveUDPAddr("udp", dstAddr)

	ln, err := net.DialUDP("udp", udpSrcAddr, udpDstAddr)
	if err != nil {
		return err
	}
	ups.s = ln
	ups.mode = 1
	ups.stat = Connecting

	return nil
}

func (ups *udpSocket) recv() {
	defer ups.netWait.Done()
	for {
		if ups.stat != Connecting && ups.stat != Connected {
			goto read_end
		}

		var inBuf []byte
		n, addr, err := ups.s.ReadFrom(inBuf)
		if err != nil {
			goto read_error
		}

		if ups.stat != Connected {
			continue
		}

		ups.i.ReadBytes += uint64(n)
		ups.i.ReadLastTime = timer.Now()

		udpAddr, _ := net.ResolveUDPAddr(addr.Network(), addr.String())

		actor.DefaultSchedulerContext.Send(ups.operator, &NetChunk{Data: inBuf, Addr: udpAddr.IP, Port: uint16(udpAddr.Port)})
	}
read_error:
	ups.stat = Closing
	ups.s.Close()
read_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	ups.so.l.Lock()
	closeHandle = ups.h
	closeOperator = ups.operator
	close(ups.out)
	//-----等待写协程结束------
	for {
		if atomic.CompareAndSwapInt32(&ups.outStat, 1, 1) {
			break
		}
	}

	ups.so.s = nil
	ups.so.b = resIdle
	ups.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, NetClose{Handle: closeHandle})
}

func (ups *udpSocket) write() {
	for {
		if ups.stat != Connecting && ups.stat != Connected {
			goto write_end
		}

		select {
		case msg := <-ups.out:
			if ups.stat != Connecting && ups.stat != Connected {
				goto write_end
			}

			udpAddr, _ := net.ResolveUDPAddr("udp", fmt.Sprint(msg.Addr.String(), ":", msg.Port))
			n, err := ups.s.WriteToUDP(msg.Data, udpAddr)
			if err != nil {
				//?
				goto write_error
			}

			ups.i.WriteBytes += uint64(n)
			ups.i.WriteLastTime = timer.Now()
		}
	}
write_error:
	ups.stat = Closing
write_end:
	ups.netWait.Done()
	ups.outStat = 1
}

func (ups *udpSocket) push(data *NetChunk, n int) error {
	//? 是否可以优化
	ups.out <- data
	return nil
}

func (ups *udpSocket) setKeepAive(keep uint64) {

}

func (ups *udpSocket) getKeepAive() uint64 {
	return 0
}

func (ups *udpSocket) getLastActivedTime() uint64 {
	return ups.i.ReadLastTime
}

func (ups *udpSocket) getStat() int32 {
	return ups.stat
}

func (ups *udpSocket) setConnected() bool {
	return atomic.CompareAndSwapInt32(&ups.stat, Connecting, Connected)
}

func (ups *udpSocket) close(lck *util.ReSpinLock) {
	if lck != nil {
		lck.Lock()
	}

	if ups.stat != Closing {
		ups.stat = Closing
		ups.s.Close()
	}

	if lck != nil {
		lck.Unlock()
	}
}

func (ups *udpSocket) closewait() {
	ups.netWait.Wait()
}

func (ups *udpSocket) getProto() string {
	return ProtoUDP
}

func (ups *udpSocket) getType() int {
	if ups.mode == 0 {
		return CListen
	}
	return CClient
}
