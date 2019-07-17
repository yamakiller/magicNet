package network

import (
	"errors"
	"magicNet/engine/actor"
	"magicNet/timer"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type sServer struct {
	handle     int32
	s          interface{}
	maker      makeConn
	operator   *actor.PID
	conns      sync.Map
	outChanMax int
	netWait    sync.WaitGroup
	stat       int32
	isShutdown bool
}

type makeConn func(handle int32, s interface{}, operator *actor.PID, so *slot, now uint64, stat int32) ISocket

func (ss *sServer) serve(ln net.Listener) {

}

func (ss *sServer) accept(conn interface{},
	network string,
	address string) error {

	handle, so := operGrap()
	if handle == -1 || so == nil {
		return errors.New("lack of socket resources")
	}

	now := timer.Now()
	so.l.Lock()
	if ss.isShutdown {
		so.b = resIdle
		so.l.Unlock()
		return errors.New("server closed")
	}

	c := ss.maker(handle, conn, ss.operator, so, now, Connecting)
	so.s = c

	go c.recv()
	go c.write()

	so.b = resAssigned
	so.l.Unlock()

	ss.conns.Store(handle, int32(1))

	ip, _ := net.ResolveIPAddr(network, address)
	actor.DefaultSchedulerContext.Send(ss.operator, &NetAccept{Handle: handle, Addr: ip.IP.To16()})

	return nil
}

func (ss *sServer) keeploop() {
	defer ss.netWait.Done()
	for {
		if ss.isShutdown {
			break
		}

		time.Sleep(time.Second * 1)
		now := timer.Now()
		ss.conns.Range(func(handle interface{}, v interface{}) bool {
			so := operGet(handle.(int32))
			if so == nil {
				return true
			}

			if so.b == resIdle {
				return true
			}

			so.l.Lock()
			defer so.l.Unlock()
			if so.b == resIdle || so.b == resOccupy || so.s == nil {
				return true
			}

			// 维护KeepAlive
			if so.s.getKeepAive() == 0 {
				return true
			}

			if (now - so.s.getLastActivedTime()) > so.s.getKeepAive() {
				so.s.close(nil)
			}
			return true
		})
	}

	//------------------关闭所有连接-----------------------------
	ss.conns.Range(func(handle interface{}, v interface{}) bool {
		so := operGet(handle.(int32))
		if so.b == resIdle {
			return true
		}

		so.l.Lock()
		if so.b == resIdle || so.b == resOccupy || so.s == nil {
			so.l.Unlock()
			return true
		}

		conn := so.s
		conn.close(nil)
		so.l.Unlock()
		conn.closewait()

		return true
	})
}

func (ss *sServer) recv() {

}

func (ss *sServer) write() {

}

func (ss *sServer) setKeepAive(keep uint64) {

}

func (ss *sServer) getKeepAive() uint64 {
	return 0
}

func (ss *sServer) getLastActivedTime() uint64 {
	return 0
}

func (ss *sServer) connect(operator *actor.PID, addr string) error {
	return nil
}

func (ss *sServer) getStat() int32 {
	return ss.stat
}

func (ss *sServer) setConnected() bool {
	return atomic.CompareAndSwapInt32(&ss.stat, Connecting, Connected)
}

func (ss *sServer) closewait() {
	ss.netWait.Wait()
}

func (ss *sServer) push(data []byte, n int) {

}
