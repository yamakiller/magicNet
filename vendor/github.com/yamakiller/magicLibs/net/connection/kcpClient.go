package connection

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/mmath"
	"github.com/yamakiller/magicLibs/net/middle"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/mgokcp/mkcp"
)

//KCPSeria KCP序列化反序列化接口
type KCPSeria interface {
	UnSeria([]byte) (interface{}, error)
	Seria(interface{}, *mkcp.KCP) (int, error)
}

type recvData struct {
	_data []byte
	_len  int
}

//KCPClient KCP(UDP)协议客户端
type KCPClient struct {
	WriteWaitQueue int
	ReadWaitQueue  int
	RecvWndSize    int32
	SendWndSize    int32
	NoDelay        int32
	Interval       int32
	Resend         int32
	Nc             int32
	RxMinRto       int32
	FastResend     int32
	S              KCPSeria
	E              Exception
	Mtu            int
	Middleware     middle.KCMiddleware
	Allocator      func(int) []byte
	Releaser       func([]byte)

	_c    *net.UDPConn
	_id   uint32
	_kcp  *mkcp.KCP
	_sync sync.Mutex
	_addr *net.UDPAddr

	_cancel  context.CancelFunc
	_ctx     context.Context
	_sdQueue chan interface{}
	_rdQueue chan *recvData

	_buffer     []byte
	_wTotal     int
	_rTotal     int
	_lastActive int64

	_wg sync.WaitGroup
}

//Connect 连接服务器
func (slf *KCPClient) Connect(addr string, timeout time.Duration) error {

	udpAddr, err := net.ResolveUDPAddr("", addr)
	if err != nil {
		return err
	}

	if slf.Mtu == 0 {
		slf.Mtu = 1400
	}

	c, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return err
	}
	slf._c = c

	if slf.Middleware != nil {
		conv, err := slf.Middleware.Subscribe(c, udpAddr, timeout)
		if err != nil {
			return err
		}
		slf._id = conv.(uint32)
		c.SetReadDeadline(time.Time{})
	}

	slf._addr = udpAddr
	slf._sdQueue = make(chan interface{}, slf.WriteWaitQueue)
	slf._rdQueue = make(chan *recvData, slf.ReadWaitQueue)

	slf._buffer = make([]byte, mmath.Align(uint32(slf.Mtu), 4))
	slf._kcp = mkcp.New(slf._id, slf)
	slf._kcp.WithOutput(output)
	slf._kcp.WndSize(slf.SendWndSize, slf.RecvWndSize)
	slf._kcp.NoDelay(slf.NoDelay, slf.Interval, slf.Resend, slf.Nc)
	slf._kcp.SetMTU(int32(slf.Mtu))
	if slf.RxMinRto > 0 {
		slf._kcp.SetRxMinRto(slf.RxMinRto)
	}
	if slf.FastResend > 0 {
		slf._kcp.SetFastResend(slf.FastResend)
	}

	slf._wg.Add(2)
	ctx, cancel := context.WithCancel(context.Background())
	slf._cancel = cancel
	slf._ctx = ctx
	go slf.writeServe()
	go slf.readServe()

	return nil
}

func (slf *KCPClient) writeServe() {
	defer func() {
		slf._wg.Done()
	}()

	var current int64
	for {
	active:
		current = util.Timestamp()
		select {
		case <-slf._ctx.Done():
			goto exit
		case <-time.After(time.Duration(slf.Interval) * time.Millisecond):
		case msg, ok := <-slf._sdQueue:
			if !ok {
				goto exit
			}

			slf._sync.Lock()
			n, err := slf.S.Seria(msg, slf._kcp)
			slf._sync.Unlock()
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}

				goto active
			}

			slf._wTotal += n
		}

		slf._sync.Lock()
		slf._kcp.Update(uint32(current & 0xFFFFFFFF))
		slf._sync.Unlock()
	}
exit:
}

func (slf *KCPClient) readServe() {
	defer func() {
		slf._wg.Done()
	}()

	for {
	active:
		select {
		case <-slf._ctx.Done():
			goto exit
		default:
			n, _, err := slf._c.ReadFromUDP(slf._buffer)
			if err != nil {
				if slf.E != nil {
					slf.E.Error(err)
				}
				goto active
			}

			slf._sync.Lock()
			slf._kcp.Input(slf._buffer, int32(n))
			slf._sync.Unlock()
			slf._rTotal += n
			slf._lastActive = util.Timestamp()

			for {
				slf._sync.Lock()
				n = int(slf._kcp.Recv(slf._buffer, int32(len(slf._buffer))))
				slf._sync.Unlock()
				if n < 0 {
					break
				}
				//需要修改池化
				var tmpBuf []byte
				if slf.Allocator != nil {
					tmpBuf = slf.Allocator(n)
				} else {
					tmpBuf = make([]byte, n)
				}

				copy(tmpBuf, slf._buffer[:n])
				slf._rdQueue <- &recvData{_data: tmpBuf, _len: n}
			}
		}
	}
exit:
}

//Parse 解析数据, 需要修改
func (slf *KCPClient) Parse() (interface{}, error) {
	slf._wg.Add(1)
	defer slf._wg.Done()

	select {
	case <-slf._ctx.Done():
		return nil, errors.New("closed")
	case d, ok := <-slf._rdQueue:
		if !ok {
			return nil, errors.New("closed")
		}

		defer func() {
			if slf.Releaser != nil {
				slf.Releaser(d._data)
			}
		}()

		msg, err := slf.S.UnSeria(d._data[:d._len])
		if err != nil {
			return nil, err
		}

		return msg, nil
	}
}

//SendTo 发送数据
func (slf *KCPClient) SendTo(msg interface{}) error {
	slf._wg.Add(1)
	defer slf._wg.Done()

	select {
	case <-slf._ctx.Done():
		return errors.New("closed")
	default:
	}

	slf._sdQueue <- msg
	return nil
}

//Close 关闭
func (slf *KCPClient) Close() error {
	select {
	case <-slf._ctx.Done():
		slf._wg.Wait()
		return errors.New("closed")
	default:
	}

	if slf._cancel != nil {
		slf._cancel()
	}

	err := slf._c.Close()
	slf._wg.Wait()

	if slf._rdQueue != nil {
		close(slf._rdQueue)
		for d := range slf._rdQueue {
			if d == nil {
				break
			}
			if slf.Releaser != nil {
				slf.Releaser(d._data)
			}

		}
		slf._rdQueue = nil
	}

	if slf._sdQueue != nil {

		close(slf._sdQueue)
		slf._sdQueue = nil
	}

	if slf._kcp != nil {
		mkcp.Free(slf._kcp)
		slf._kcp = nil
	}

	return err
}

func output(buff []byte, user interface{}) int32 {
	client := user.(*KCPClient)
	client._c.Write(buff)
	return 0
}
