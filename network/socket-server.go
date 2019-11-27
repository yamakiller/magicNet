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
	_h          int32
	_s          interface{}
	_maker      makeConn
	_operator   *actor.PID
	_conns      sync.Map
	_outChanMax int
	_netWait    sync.WaitGroup
	_stat       int32
	_isShutdown bool
}

type makeConn func(handle int32, s interface{}, operator *actor.PID, so *slot, now uint64, stat int32) ISocket

func (slf *sServer) serve(ln net.Listener) {

}

func (slf *sServer) accept(conn interface{},
	network string,
	address string) error {

	handle, so := operGrap()
	if handle == -1 || so == nil {
		return errors.New("lack of socket resources")
	}

	now := timer.Now()
	so.l.Lock()
	if slf._isShutdown {
		so.b = resIdle
		so.l.Unlock()
		return errors.New("server closed")
	}

	c := slf._maker(handle, conn, slf._operator, so, now, Connecting)
	so.s = c

	go c.recv()
	go c.write()

	so.b = resAssigned
	so.l.Unlock()

	slf._conns.Store(handle, int32(1))

	addr, _ := net.ResolveTCPAddr(network, address)

	actor.DefaultSchedulerContext.Send(slf._operator, &NetAccept{Handle: handle, Addr: addr.IP.To16(), Port: addr.Port})

	return nil
}

func (slf *sServer) recv() {

}

func (slf *sServer) write() {

}

func (slf *sServer) setKeepAive(keep uint64) {

}

func (slf *sServer) getKeepAive() uint64 {
	return 0
}

func (slf *sServer) getLastActivedTime() uint64 {
	return 0
}

func (slf *sServer) connect(operator *actor.PID, addr string) error {
	return nil
}

func (slf *sServer) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (slf *sServer) getStat() int32 {
	return slf._stat
}

func (slf *sServer) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf._stat, Connecting, Connected)
}

func (slf *sServer) closewait() {
	slf._netWait.Wait()
}

func (slf *sServer) push(data *NetChunk, n int) error {
	return nil
}
