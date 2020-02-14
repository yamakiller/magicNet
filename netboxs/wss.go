package netboxs

import (
	"errors"
	"strconv"
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
	boxs.Box
	_max    int32
	_cur    int32
	_borker *borker.WSSBorker
	_conns  *table.HashTable
	_sync   sync.Mutex
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
//addr@wspath@readbuffersize@writebuffersize@handshaketimeout(second)
func (slf *WSSBox) ListenAndServe(addr string) error {
	slf.Box.StartedWait()
	wsPath := ""
	as := strings.Split(addr, "@")
	rBufferSize := 2048
	wBufferSize := 2048
	hTimeOut := time.Second
	if len(as) >= 2 {
		addr = as[0]
		wsPath = as[1]
		if len(as) >= 3 {
			rBufferSize, _ = strconv.Atoi(as[2])
		}

		if len(as) >= 4 {
			wBufferSize, _ = strconv.Atoi(as[2])
		}

		if len(as) >= 5 {
			tmpTimeOut, _ := strconv.Atoi(as[2])
			hTimeOut = time.Duration(tmpTimeOut) * time.Second
		}
	}

	slf._borker = &borker.WSSBorker{
		WSPath:           wsPath,
		Spawn:            slf.handleConnect,
		ReadBufferSize:   rBufferSize,
		WriteBufferSize:  wBufferSize,
		HandshakeTimeout: hTimeOut,
	}

	slf._conns = &table.HashTable{
		Mask: 0xFFFFFFF,
		Max:  uint32(float64(slf._max) * 1.2),
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

//CloseConn 关闭一个连接
func (slf *WSSBox) CloseConn(socket int32) error {
	slf._sync.Lock()
	defer slf._sync.Unlock()
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		return errors.New("not found socket")
	}
	cc := c.(*_WBoxConn)
	cc._state = stateClosed
	cc._closed <- true
	err := cc._io.Close()
	return err
}

//CloseConnWait 关闭一个连接并等待连接退出
func (slf *WSSBox) CloseConnWait(socket int32) error {
	slf._sync.Lock()
	c := slf._conns.Get(uint32(socket))
	if c == nil {
		slf._sync.Unlock()
		return errors.New("not found socket")
	}

	cc := c.(*_WBoxConn)
	cc._state = stateClosed
	cc._closed <- true
	err := cc._io.Close()
	cc._wg.Wait()

	return err
}

//GetValues Returns all socket
func (slf *WSSBox) GetValues() []int32 {
	slf._sync.Lock()
	defer slf._sync.Unlock()

	cns := slf._conns.GetValues()
	res := make([]int32, len(cns))
	for k, c := range cns {
		res[k] = c.(*_TBoxConn)._cn.Socket()
	}

	return res
}

func (slf *WSSBox) handleCloseAll() {
	cs := slf.GetValues()
	for {
		for _, socket := range cs {
			slf.CloseConnWait(socket)
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

	cc := &_WBoxConn{
		_io:     c,
		_closed: make(chan bool),
		_cn:     slf._pools.Get(),
		_state:  stateInit,
	}

	slf._sync.Lock()
	s, err := slf._conns.Push(cc)
	if err != nil {
		atomic.AddInt32(&slf._cur, -1)
		slf._sync.Unlock()
		return err
	}
	slf._sync.Unlock()
	socket := int32(s)

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
					Add(time.Duration(float64(cc._cn.Keepalive()) * 1.5)))
			}

			msg, err := cc._cn.Parse()
			if cc._state == stateClosed {
				err = errors.New("error disconnect")
			}

			if msg != nil {
				slf.Box.GetPID().Post(&netmsgs.Message{Sock: socket,
					Data: msg})
			}

			if err != nil {
				slf.CloseConn(socket)
				slf.Box.GetPID().Post(&netmsgs.Closed{Sock: socket})
			}
		}
	}()

	go func() {
		defer func() {
			cc._cn.Close()
			cc._wg.Done()
		}()

		for {
			select {
			case <-cc._closed:
				goto exit
			case <-cc._kicker.C:
				if cc._cn.Keepalive() > 0 {
					cc._cn.Ping()
					cc._kicker.Reset(cc._cn.Keepalive())
				}
			case msg := <-cc._cn.Pop():
				state := cc._state
				if state == stateConnected || state == stateConnecting {
					if err := cc._cn.Write(msg); err != nil {
						slf.Box.GetPID().Post(&netmsgs.Error{Sock: socket, Err: err})
					}

					if cc._kicker != nil && cc._cn.Keepalive() > 0 {
						cc._kicker.Reset(cc._cn.Keepalive())
					}
					cc._activity = time.Now()
				}
			}
		}
	exit:
	}()

	slf.Box.GetPID().Post(&netmsgs.Accept{Sock: socket})

	return nil
}

type _WBoxConn struct {
	_io       *listener.WSSConn
	_cn       Connect
	_wg       sync.WaitGroup
	_state    state
	_closed   chan bool
	_kicker   *time.Timer
	_activity time.Time
}
