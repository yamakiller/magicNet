package netboxs

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/net/borker"
	"github.com/yamakiller/magicLibs/net/listener"
	"github.com/yamakiller/magicLibs/st/table"
	"github.com/yamakiller/magicNet/netmsgs"
)

//KCPBox kcp network box
type KCPBox struct {
	RecvWndSize   int32
	SendWndSize   int32
	RecvQueueSize int32
	NoDelay       int32
	Interval      int32
	Resend        int32
	Nc            int32
	RxMinRto      int32
	FastResend    int32
	Mtu           int

	boxs.Box
	_max    int32
	_cur    int32
	_borker *borker.KCPBorker
	_conns  *table.HashTable2
	_sync   sync.Mutex
	_closed bool
	_pools  Pool
}

//WithPool setting connection pools
func (slf *KCPBox) WithPool(pool Pool) {
	slf._pools = pool
}

//WithMax setting connection max of number
func (slf *KCPBox) WithMax(max int32) {
	slf._max = max
}

//ListenAndServe 启动监听服务
func (slf *KCPBox) ListenAndServe(addr string) error {
	slf.Box.StartedWait()
	slf._borker = &borker.KCPBorker{
		RecvWndSize:   slf.RecvWndSize,
		SendWndSize:   slf.SendWndSize,
		RecvQueueSize: slf.RecvQueueSize,
		NoDelay:       slf.NoDelay,
		Interval:      slf.Interval,
		Resend:        slf.Resend,
		Nc:            slf.Nc,
		RxMinRto:      slf.RxMinRto,
		FastResend:    slf.FastResend,
		Mtu:           slf.Mtu,
		Spawn:         slf.handleConnect,
	}

	slf._conns = &table.HashTable2{
		Mask: 0xFFFFFFF,
		Max:  uint32(slf._max),
		Comp: func(a, b interface{}) int {
			ca := a.(*_KBoxConn)
			cb := b.(uint32)
			if ca._cn.Socket() == int32(cb) {
				return 0
			}
			return -1
		},
		GetKey: func(a interface{}) uint32 {
			return uint32(a.(*_KBoxConn)._cn.Socket())
		},
	}

	slf._conns.Initial()

	if err := slf._borker.ListenAndServe(addr); err != nil {
		return err
	}

	return nil
}

//OpenTo setting connection state connected
func (slf *KCPBox) OpenTo(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	cc := c.(*_WBoxConn)
	cc._state = stateConnected

	return nil
}

//SendTo 发送数据给连接
func (slf *KCPBox) SendTo(socket int32, msg interface{}) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	cc := c.(*_TBoxConn)
	if cc._state == stateClosed {
		return errors.New("connection closed")
	}

	select {
	case <-cc._closed:
	default:
	}

	return cc._cn.Push(msg)
}

//CloseTo 关闭一个连接
func (slf *KCPBox) CloseTo(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}
	slf._conns.Remove(uint32(socket))
	cc := c.(*_KBoxConn)
	cc._state = stateClosed
	select {
	case <-cc._closed:
	default:
		close(cc._closed)
	}
	err := cc._io.Close()
	return err
}

//CloseToWait 关闭一个连接并等待连接退出
func (slf *KCPBox) CloseToWait(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}
	slf._conns.Remove(uint32(socket))

	cc := c.(*_KBoxConn)
	cc._state = stateClosed
	select {
	case <-cc._closed:
	default:
		close(cc._closed)
	}

	err := cc._io.Close()
	cc._wg.Wait()

	return err
}

//GetValues Returns all socket
func (slf *KCPBox) GetValues() []int32 {
	cns := slf._conns.GetValues()
	res := make([]int32, len(cns))
	for k, c := range cns {
		res[k] = c.(*_KBoxConn)._cn.Socket()
	}

	return res
}

//Shutdown 关闭服务
func (slf *KCPBox) Shutdown() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

//ShutdownWait 关闭服务并等待结束
func (slf *KCPBox) ShutdownWait() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

func (slf *KCPBox) handleCloseAll() {
	cs := slf.GetValues()
	for {
		for _, socket := range cs {
			slf.CloseToWait(socket)
		}

		cs = slf.GetValues()
		if len(cs) <= 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (slf *KCPBox) handleConnect(c *listener.KCPConn) error {
	if slf._closed {
		return errors.New("listener closed")
	}

	newSz := atomic.AddInt32(&slf._cur, 1)
	if newSz > slf._max {
		atomic.AddInt32(&slf._cur, -1)
		return errors.New("connection is full")
	}

	cc := &_KBoxConn{
		_io:       c,
		_closed:   make(chan bool, 1),
		_cn:       slf._pools.Get(),
		_activity: time.Now(),
		_state:    stateInit,
	}

	s, err := slf._conns.Push(cc)
	if err != nil {
		atomic.AddInt32(&slf._cur, -1)
		return err
	}

	socket := int32(s)
	cc._cn.WithIO(c)
	cc._cn.WithSocket(socket)
	if cc._cn.Keepalive() > 0 {
		cc._kicker = time.NewTimer(cc._cn.Keepalive())
	}
	cc._state = stateAccepted
	cc._wg.Add(1)
	//数据接收
	go func() {
		for {
			if cc._cn.Keepalive() > 0 {
				cc._io.SetReadDeadline(cc._activity.
					Add(time.Duration(float64(cc._cn.Keepalive()) * 2.0)))
			}

			msg, err := cc._cn.UnSeria()
			if cc._state == stateClosed {
				err = errors.New("error disconnect")
			}

			if msg != nil {
				slf.Box.GetPID().Post(&netmsgs.Message{Sock: socket,
					Data: msg})
			}

			if err != nil {
				slf.CloseTo(socket)
				slf.Box.GetPID().Post(&netmsgs.Closed{Sock: socket})
				break
			}
		}
	}()

	//数据发送
	go func() {
		defer func() {
			cc._cn.Close()
			slf._pools.Put(cc._cn)
			cc._wg.Done()
		}()

		for {
		active:
			if cc._kicker != nil {

				select {
				case <-cc._closed:
					goto exit
				case <-cc._kicker.C:
					if cc._cn.Keepalive() > 0 {
						cc._cn.Ping()
						cc._kicker.Reset(cc._cn.Keepalive())
					}
				case msg, ok := <-cc._cn.Pop():
					if !ok {
						goto active
					}

					if err := cc._cn.Seria(msg); err != nil && cc._state != stateClosed {
						slf.Box.GetPID().Post(&netmsgs.Error{Sock: socket, Err: err})
					}
				}
			} else {
				select {
				case <-cc._closed:
					goto exit
				case msg := <-cc._cn.Pop():
					if err := cc._cn.Seria(msg); err != nil && cc._state != stateClosed {
						slf.Box.GetPID().Post(&netmsgs.Error{Sock: socket, Err: err})
					}
				}
			}
		}
	exit:
	}()

	return nil
}

type _KBoxConn struct {
	_io       *listener.KCPConn
	_cn       Connect
	_wg       sync.WaitGroup
	_state    state
	_closed   chan bool
	_kicker   *time.Timer
	_activity time.Time
}

//BKCPConn KCP base connection
type BKCPConn struct {
	WriteQueueSize int
	_readWriter    *listener.KCPConn
	_sock          int32
	_queue         chan interface{}
}

//Socket Returns socket
func (slf *BKCPConn) Socket() int32 {
	return slf._sock
}

//WithSocket setting socket
func (slf *BKCPConn) WithSocket(sock int32) {
	slf._sock = sock
}

//WithIO setting io interface
func (slf *BKCPConn) WithIO(c interface{}) {
	slf._readWriter = c.(*listener.KCPConn)
	slf._queue = make(chan interface{}, slf.WriteQueueSize)
}

//Reader Returns reader buffer
func (slf *BKCPConn) Reader() *listener.KCPConn {
	return slf._readWriter
}

//Writer Returns writer buffer
func (slf *BKCPConn) Writer() *listener.KCPConn {
	return slf._readWriter
}

//Push 插入发送数据
func (slf *BKCPConn) Push(msg interface{}) error {
	slf._queue <- msg
	return nil
}

//Pop 弹出需要发送的数据
func (slf *BKCPConn) Pop() chan interface{} {
	return slf._queue
}

//Close 释放连接资源
func (slf *BKCPConn) Close() error {
	if slf._queue != nil {
		close(slf._queue)
	}
	return nil
}
