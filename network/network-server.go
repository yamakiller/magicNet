package network

import (
	"errors"
	"magicNet/engine/actor"
	"magicNet/engine/util"
	"sync/atomic"
)

const (
	maxSocket = 65535
)

const (
	resIdle     = 0
	resOccupy   = 1
	resAssigned = 2
)

const (
	// ErrSocketResources : 套接字资源不足
	ErrSocketResources = "lack of socket resources"
)

// OperWSListen : 开启一个websocket 监听
func OperWSListen(operator *actor.PID, addr string) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	// ? 是否考虑被 resIdle 得处理
	/*s.s = &wslisten{handle: h}
	s.b = resAssigned
	lstn, _ := s.s.(*wslisten)
	if err := lstn.listen(operator, addr); err != nil {
		s.b = resIdle
		s.s = nil
		s.l.Unlock()
		return -1, err
	}*/
	s.l.Unlock()

	return h, nil
}

// OperWSConnect : 创建一个websocket 客户端连接
func OperWSConnect(operator *actor.PID, addr string) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	// TODO : 连接Client
	s.l.Unlock()

	return h, nil
}

// OperClose : 关闭一个Socket
func OperClose(handle int32) {
	s := operGet(handle)
	if s == nil {
		return
	}

	s.l.Lock()
	if s.b == resIdle || s.b == resOccupy {
		s.b = resIdle
		s.l.Unlock()
		return
	}

	s.s.close(nil)
	s.l.Unlock()
	s.s.closewait()
}

func operGet(handle int32) *slot {
	return defaultNServer.get(handle)
}

func operGrap() (int32, *slot) {
	return defaultNServer.grap()
}

func operForeach(f func(s *slot)) {
	defaultNServer.foreach(f)
}

type slot struct {
	b int32
	s ISocket
	l util.SpinLock
}

// NetServer : 网络服务器管理器
type NetServer struct {
	ss []slot
	fi int32
	//sl util.SpinLock
}

var (
	defaultNServer = NetServer{ss: make([]slot, maxSocket)}
)

func (ns *NetServer) get(handle int32) *slot {
	return &ns.ss[ns.hash(handle)]
}

func (ns *NetServer) grap() (int32, *slot) {

	for i := 0; i < maxSocket; i++ {
		ns.fi++
		handle := ns.fi
		if handle < 0 {
			ns.fi &= 0x7FFFFFFF
			handle = ns.fi
		}
		s := &ns.ss[ns.hash(handle)]
		if s.b == resIdle {
			if atomic.CompareAndSwapInt32(&s.b, resIdle, resOccupy) {
				return handle, s
			}
			i--
		}
	}

	return -1, nil
}

func (ns *NetServer) foreach(f func(so *slot)) {
	for i := 0; i < maxSocket; i++ {
		s := &ns.ss[i]
		if s.b == resIdle {
			continue
		}

		f(s)
	}
}

func (ns *NetServer) hash(handle int32) int32 {
	return handle & (maxSocket - 1)
}
