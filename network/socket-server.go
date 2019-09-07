package network

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/timer"
)

type sServer struct {
	h          int32
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

	tpc, _ := net.ResolveTCPAddr(network, address)
	actor.DefaultSchedulerContext.Send(ss.operator, &NetAccept{Handle: handle, Addr: tpc.IP.To16(), Port: tpc.Port})

	return nil
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

func (ss *sServer) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
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

func (ss *sServer) push(data *NetChunk, n int) error {
	return nil
}
