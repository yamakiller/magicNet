package netboxs

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/net/borker"
	"github.com/yamakiller/magicLibs/net/listener"
	"github.com/yamakiller/magicLibs/st/table"
	"github.com/yamakiller/magicNet/netmsgs"
)

//WSSBox websocket network box
type WSSBox struct {
	ReadBufferSize   int
	WriteBufferSize  int
	Handshaketimeout int
	T                int //消息类型TextMessage ...

	boxs.Box
	_max    int32
	_cur    int32
	_borker *borker.WSSBorker
	_conns  *table.HashTable2

	_closed bool
	_pools  Pool
}

//WithPool setting connection pools
func (slf *WSSBox) WithPool(pool Pool) {
	slf._pools = pool
}

//WithMax setting connection max of number
func (slf *WSSBox) WithMax(max int32) {
	slf._max = max
}

//ListenAndServe 启动监听服务
//addr@wspath
func (slf *WSSBox) ListenAndServe(addr string) error {
	slf.Box.StartedWait()
	wsPath := ""
	as := strings.Split(addr, "@")
	hTimeOut := time.Second
	if slf.Handshaketimeout > 0 {
		hTimeOut = slf._borker.HandshakeTimeout * time.Second
	}

	if len(as) >= 2 {
		addr = as[0]
		wsPath = as[1]
	}

	slf._borker = &borker.WSSBorker{
		WSPath:           wsPath,
		Spawn:            slf.handleConnect,
		ReadBufferSize:   slf.ReadBufferSize,
		WriteBufferSize:  slf.WriteBufferSize,
		HandshakeTimeout: hTimeOut,
	}

	slf._conns = &table.HashTable2{
		Mask: 0xFFFFFFF,
		Max:  uint32(slf._max),
		Comp: func(a, b interface{}) int {
			ca := a.(*_WBoxConn)
			cb := b.(uint32)
			if ca._cn.Socket() == int32(cb) {
				return 0
			}
			return -1
		},
		GetKey: func(a interface{}) uint32 {
			return uint32(a.(*_WBoxConn)._cn.Socket())
		},
	}
	slf._conns.Initial()

	if err := slf._borker.ListenAndServe(addr); err != nil {
		return err
	}

	return nil
}

//Shutdown 关闭服务
func (slf *WSSBox) Shutdown() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

//ShutdownWait 关闭服务并等待结束
func (slf *WSSBox) ShutdownWait() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

//SendTo 发送数据给连接
func (slf *WSSBox) SendTo(socket int32, msg interface{}) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	cc := c.(*_WBoxConn)
	if cc._state == stateClosed {
		return errors.New("connection closed")
	}

	cc._swg.Add(1)
	defer cc._swg.Done()

	select {
	case <-cc._ctx.Done():
	default:
	}

	return cc._cn.Push(msg)
}

//CloseTo 关闭一个连接
func (slf *WSSBox) CloseTo(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	slf._conns.Remove(uint32(socket))
	atomic.AddInt32(&slf._cur, -1)
	cc := c.(*_WBoxConn)
	cc._state = stateClosed
	cc._cancel()
	err := cc._io.Close()
	return err
}

//CloseToWait 关闭一个连接并等待连接退出
func (slf *WSSBox) CloseToWait(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	if slf._conns.Remove(uint32(socket)) {
		atomic.AddInt32(&slf._cur, -1)
	}
	cc := c.(*_WBoxConn)
	cc._state = stateClosed
	cc._cancel()
	err := cc._io.Close()
	cc._wg.Wait()

	return err
}

//GetConnect Return connection
func (slf *WSSBox) GetConnect(socket int32) (interface{}, error) {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return nil, errors.New("not found connection")
	}
	return c.(*_WBoxConn)._cn, nil
}

//GetValues Returns all socket
func (slf *WSSBox) GetValues() []int32 {
	cns := slf._conns.GetValues()
	res := make([]int32, len(cns))
	for k, c := range cns {
		res[k] = c.(*_WBoxConn)._cn.Socket()
	}

	return res
}

func (slf *WSSBox) handleCloseAll() {
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

func (slf *WSSBox) handleConnect(c *listener.WSSConn) error {
	if slf._closed {
		return errors.New("listener closed")
	}

	newSz := atomic.AddInt32(&slf._cur, 1)
	if newSz > slf._max {
		atomic.AddInt32(&slf._cur, -1)
		return errors.New("connection is full")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cc := &_WBoxConn{
		_io:     c,
		_cancel: cancel,
		_ctx:    ctx,
		_cn:     slf._pools.Get(),
		_state:  stateInit,
	}

	s, err := slf._conns.Push(cc)
	if err != nil {
		atomic.AddInt32(&slf._cur, -1)
		return err
	}

	socket := int32(s)
	socketAddr := c.RemoteAddr()

	cc._cn.WithIO(c)
	cc._cn.WithSocket(socket)
	if cc._cn.Keepalive() > 0 {
		cc._kicker = time.NewTimer(cc._cn.Keepalive())
	}
	cc._state = stateAccepted

	cc._wg.Add(1)
	go func() {
		for {
			if cc._cn.Keepalive() > 0 {
				cc._io.Conn.SetReadDeadline(cc._activity.
					Add(time.Duration(float64(cc._cn.Keepalive()) * 2)))
			}

			msg, err := cc._cn.UnSeria()
			if cc._state == stateClosed {
				err = errors.New("error disconnect")
			}

			if msg != nil {
				cc._activity = time.Now()
				slf.Box.GetPID().Post(&netmsgs.Message{Sock: socket,
					Data: msg})
				state := cc._state
				if state == stateConnected || state == stateConnecting {

					if cc._kicker != nil && cc._cn.Keepalive() > 0 {
						cc._kicker.Reset(cc._cn.Keepalive())
					}
					cc._activity = time.Now()
				}
			}

			if err != nil {
				slf.CloseTo(socket)
				slf.Box.GetPID().Post(&netmsgs.Closed{Sock: socket})
			}
		}
	}()

	go func() {
		defer func() {
			cc._swg.Wait()
			cc._cn.Close()
			cc._wg.Done()
		}()

		for {
		active:
			if cc._kicker != nil {
				select {
				case <-cc._ctx.Done():
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
				case <-cc._ctx.Done():
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

	slf.Box.GetPID().Post(&netmsgs.Accept{Sock: socket, Addr: socketAddr})

	return nil
}

type _WBoxConn struct {
	_io       *listener.WSSConn
	_cn       Connect
	_wg       sync.WaitGroup
	_swg      sync.WaitGroup
	_state    state
	_cancel   context.CancelFunc
	_ctx      context.Context
	_kicker   *time.Timer
	_activity time.Time
}

//BWSSConn WebSocket conn base
type BWSSConn struct {
	WriteQueueSize int

	_c     *listener.WSSConn
	_t     int
	_sock  int32
	_queue chan interface{}
}

//Socket Returns socket
func (slf *BWSSConn) Socket() int32 {
	return slf._sock
}

//WithSocket setting socket
func (slf *BWSSConn) WithSocket(sock int32) {
	slf._sock = sock
}

//WithIO 设置底层ID
func (slf *BWSSConn) WithIO(c interface{}) {
	slf._c = c.(*listener.WSSConn)
	slf._queue = make(chan interface{}, slf.WriteQueueSize)
}

//Reader Returns reader buffer
func (slf *BWSSConn) Reader() (io.Reader, error) {
	t, r, err := slf._c.NextReader()
	if err != nil {
		return nil, err
	}

	if t != slf._t {
		return nil, errors.New("data format does not match")
	}

	return r, nil
}

//Writer Returns writer buffer
func (slf *BWSSConn) Writer() (io.WriteCloser, error) {
	w, err := slf._c.NextWriter(slf._t)
	if err != nil {
		return nil, err
	}

	return w, nil
}

//Push 插入发送数据
func (slf *BWSSConn) Push(msg interface{}) error {
	slf._queue <- msg
	return nil
}

//Pop 弹出需要发送的数据
func (slf *BWSSConn) Pop() chan interface{} {
	return slf._queue
}

//Close 释放连接资源
func (slf *BWSSConn) Close() error {
	if slf._queue != nil {
		close(slf._queue)
	}
	return nil
}
