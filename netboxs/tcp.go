package netboxs

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/net/borker"
	"github.com/yamakiller/magicLibs/st/table"
	"github.com/yamakiller/magicNet/netmsgs"
)

// TCPBox tcp network box
type TCPBox struct {
	boxs.Box
	_max    int32
	_cur    int32
	_borker *borker.TCPBorker
	_conns  *table.HashTable2

	_closed bool
	_pools  Pool
}

// WithPool setting connection pools
func (slf *TCPBox) WithPool(pool Pool) {
	slf._pools = pool
}

// WithMax setting connection max of number
func (slf *TCPBox) WithMax(max int32) {
	slf._max = max
}

// ListenAndServe 启动监听服务
func (slf *TCPBox) ListenAndServe(addr string) error {
	slf.Box.StartedWait()
	slf._borker = &borker.TCPBorker{
		Spawn: slf.handleConnect,
	}

	slf._conns = &table.HashTable2{
		Mask: 0xFFFFFFF,
		Max:  uint32(slf._max),
		Comp: func(a, b interface{}) int {
			ca := a.(*_t_connector)
			cb := b.(uint32)
			if ca._cn.Socket() == int32(cb) {
				return 0
			}
			return -1
		},
		GetKey: func(a interface{}) uint32 {
			return uint32(a.(*_t_connector)._cn.Socket())
		},
	}
	slf._conns.Initial()

	if err := slf._borker.ListenAndServe(addr); err != nil {
		return err
	}

	return nil
}

func (slf *TCPBox) ListenAndServeTls(addr string, ptls *tls.Config) error {
	slf.Box.StartedWait()
	slf._borker = &borker.TCPBorker{
		Spawn: slf.handleConnect,
	}

	slf._conns = &table.HashTable2{
		Mask: 0xFFFFFFF,
		Max:  uint32(slf._max),
		Comp: func(a, b interface{}) int {
			ca := a.(*_t_connector)
			cb := b.(uint32)
			if ca._cn.Socket() == int32(cb) {
				return 0
			}
			return -1
		},
		GetKey: func(a interface{}) uint32 {
			return uint32(a.(*_t_connector)._cn.Socket())
		},
	}
	slf._conns.Initial()

	if err := slf._borker.ListenAndServeTls(addr, ptls); err != nil {
		return err
	}

	return nil
}

// Shutdown 关闭服务
func (slf *TCPBox) Shutdown() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

// ShutdownWait 关闭服务并等待结束
func (slf *TCPBox) ShutdownWait() {
	slf._closed = true
	slf.handleCloseAll()
	slf._borker.Shutdown()
	slf.Box.ShutdownWait()
}

// OpenTo setting connection state connected
func (slf *TCPBox) OpenTo(socket interface{}) error {
	c := slf._conns.Get(uint32(socket.(int32)))
	if c == nil {
		return errors.New("not found socket")
	}

	cc := c.(*_t_connector)
	cc._state = stateConnected

	return nil
}

// SendTo send data socket
func (slf *TCPBox) SendTo(socket interface{}, msg interface{}) error {
	sock, ok := socket.(int32)
	if !ok {
		return errors.New("param(1): socket is int32")
	}
	c := slf._conns.Get(uint32(sock))
	if c == nil {
		return errors.New("not found socket")
	}
	cc := c.(*_t_connector)
	if cc._state == stateClosed {
		return errors.New("connection closed")
	}

	cc._swg.Add(1)
	defer cc._swg.Done()
	select {
	case <-cc._ctx.Done():
		return errors.New("connection closed")
	default:
	}
	return cc._cn.Push(msg)
}

// CloseTo 关闭一个连接
func (slf *TCPBox) CloseTo(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	if slf._conns.Remove(uint32(socket)) {
		atomic.AddInt32(&slf._cur, -1)
	}

	cc := c.(*_t_connector)
	cc._state = stateClosed
	cc._cancel()
	err := cc._io.Close()
	return err
}

// CloseToWait 关闭一个连接并等待连接退出
func (slf *TCPBox) CloseToWait(socket int32) error {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}

	if slf._conns.Remove(uint32(socket)) {
		atomic.AddInt32(&slf._cur, -1)
	}
	cc := c.(*_t_connector)
	cc._state = stateClosed
	cc._cancel()

	err := cc._io.Close()
	cc._wg.Wait()

	return err
}

// GetConnect Return connection
func (slf *TCPBox) GetConnect(socket int32) (interface{}, error) {
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return nil, errors.New("not found connection")
	}
	return c.(*_t_connector)._cn, nil
}

// GetValues Returns all socket
func (slf *TCPBox) GetValues() []int32 {
	cns := slf._conns.GetValues()
	res := make([]int32, len(cns))
	for k, c := range cns {
		res[k] = c.(*_t_connector)._cn.Socket()
	}

	return res
}

func (slf *TCPBox) handleCloseAll() {
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

func (slf *TCPBox) handleConnect(c net.Conn) error {
	if slf._closed {
		return errors.New("listener closed")
	}

	newSz := atomic.AddInt32(&slf._cur, 1)
	if newSz > slf._max {
		atomic.AddInt32(&slf._cur, -1)
		return errors.New("connection is full")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cc := &_t_connector{
		_io:       c,
		_cancel:   cancel,
		_ctx:      ctx,
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
				if cn, ok := cc._io.(net.Conn); ok {
					cn.SetReadDeadline(cc._activity.
						Add(time.Duration(float64(cc._cn.Keepalive()) * 2.0)))
				}
			}

			msg, err := cc._cn.UnSeria()
			if cc._state == stateClosed {
				err = errors.New("error disconnect")
			}

			if msg != nil {
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
				break
			}
		}
	}()

	go func() {
		defer func() {
			cc._swg.Wait()
			cc._cn.Close()
			slf._pools.Put(cc._cn)
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

type _t_connector struct {
	_io       io.ReadWriteCloser
	_cn       Connect
	_wg       sync.WaitGroup
	_swg      sync.WaitGroup
	_state    state
	_cancel   context.CancelFunc
	_ctx      context.Context
	_kicker   *time.Timer
	_activity time.Time
}

type DefaultTcpConnector struct {
	_s int32
	_r *bufio.Reader
	_w *bufio.Writer
	_q chan interface{}
}

func (dtc *DefaultTcpConnector) Socket() int32 {
	return dtc._s
}

func (dtc *DefaultTcpConnector) WithSocket(sock int32) {
	slf._s = sock
}

// WithIO setting io interface
func (dtc *DefaultTcpConnector) WithIO(c interface{}) {
	dtc._r = bufio.NewReaderSize(c.(io.ReadWriteCloser), 1024)
	dtc._w = bufio.NewWriterSize(c.(io.ReadWriteCloser), 1024)
	dtc._q = make(chan interface{}, 8)
}

// Reader Returns reader buffer
func (dtc *DefaultTcpConnector) Reader() *bufio.Reader {
	return dtc._r
}

// Writer Returns writer buffer
func (dtc *DefaultTcpConnector) Writer() *bufio.Writer {
	return dtc._w
}

// Push 插入发送数据
func (dtc *DefaultTcpConnector) Push(msg interface{}) error {
	dtc._q <- msg
	return nil
}

// Pop 弹出需要发送的数据
func (dtc *DefaultTcpConnector) Pop() chan interface{} {
	return dtc._q
}

func (dtc *DefaultTcpConnector) Close() error {
	if dtc._q != nil {
		close(dtc._q)
	}

	dtc._r = nil
	dtc._w = nil

	return nil
}
