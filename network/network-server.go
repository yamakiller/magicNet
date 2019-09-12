package network

import (
	"errors"
	"net"
	"sync/atomic"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/util"
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
	// ErrUnknownSocket :
	ErrUnknownSocket = "unknown socket"
	// ErrClosedSocket :
	ErrClosedSocket = "socket is closed"
	// ErrUnConnected :
	ErrUnConnected = "socket is not connected"
)

// OperWSListen : 开启一个websocket 监听
func OperWSListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	s.s = &wsServer{}
	wss, _ := s.s.(*wsServer)
	wss.h = h
	wss.operator = operator
	wss.outChanMax = outChanSize

	s.b = resAssigned
	if err := wss.listen(operator, addr); err != nil {
		s.b = resIdle
		s.s = nil
		s.l.Unlock()
		return -1, err
	}
	s.l.Unlock()

	return h, nil
}

// OperWSConnect : 创建一个websocket 客户端连接
func OperWSConnect(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	client := &wsClient{}
	client.h = h
	client.o = operator
	client.out = make(chan *NetChunk, outChanSize)
	client.so = s
	err := client.connect(operator, addr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}

	s.l.Lock()
	s.s = client
	client.w.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperTCPListen : 开启一个 tcp socket 监听
func OperTCPListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	s.s = &tcpServer{}
	tps, _ := s.s.(*tcpServer)
	tps.h = h
	tps.operator = operator
	tps.outChanMax = outChanSize

	s.b = resAssigned
	if err := tps.listen(operator, addr); err != nil {
		s.b = resIdle
		s.s = nil
		s.l.Unlock()
		return -1, err
	}
	s.l.Unlock()

	return h, nil
}

// OperTCPConnect : 创建一个tcp socket 客户端连接
func OperTCPConnect(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	client := &tcpClient{}
	client.h = h
	client.o = operator
	client.out = make(chan *NetChunk, outChanSize)
	client.quit = make(chan int)
	client.so = s
	err := client.connect(operator, addr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}

	s.l.Lock()
	s.s = client
	client.w.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperUDPListen : 开启一个 udp socket 监听
func OperUDPListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	s.s = &udpSocket{}
	ups, _ := s.s.(*udpSocket)
	ups.h = h
	ups.so = s
	ups.operator = operator
	ups.out = make(chan *NetChunk, outChanSize)
	ups.quit = make(chan int)

	s.b = resAssigned
	if err := ups.listen(operator, addr); err != nil {
		s.b = resIdle
		s.s = nil
		s.l.Unlock()
		return -1, err
	}
	s.l.Unlock()

	return h, nil
}

// OperUDPConnect : 创建一个udp socket 客户端连接
func OperUDPConnect(operator *actor.PID, srcAddr string, dstAddr string, outChanSize int) (int32, error) {

	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	client := &udpSocket{}
	client.h = h
	client.operator = operator
	client.out = make(chan *NetChunk, outChanSize)
	client.quit = make(chan int)
	client.so = s

	err := client.udpConnect(operator, srcAddr, dstAddr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}
	s.l.Lock()
	s.s = client
	client.netWait.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperRPCListen 打开RPC服务
func OperRPCListen(operator *actor.PID,
	addr string,
	reg GRpcRegister, /*(注册协议调用函数)*/
	writeBufSize int, /*(可选项 默认:32 * 1024)*/
	readBufSize int, /*(可选项 默认:32 * 1024)*/
	connectionTimeout int, /*(可选项 默认:120 sec)*/
	maxSendMessageSize int, /*(可选项 默认: math.MaxInt32)*/
	maxReceiveMessageSize int /*(可选项 默认:1024 * 1024 * 4)*/) (int32, error) {

	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, errors.New(ErrSocketResources)
	}

	s.l.Lock()
	defer s.l.Unlock()
	s.s = &grpcServer{register: reg, writeBufSize: writeBufSize,
		readBufSize:           readBufSize,
		connectionTimeout:     connectionTimeout,
		maxSendMessageSize:    maxSendMessageSize,
		maxReceiveMessageSize: maxReceiveMessageSize}

	rpcs, _ := s.s.(*grpcServer)
	rpcs.h = h
	rpcs.so = s
	rpcs.operator = operator
	if err := rpcs.listen(operator, addr); err != nil {
		s.s = nil
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}

	s.b = resAssigned

	return h, nil
}

//OperWrite : 发送数据
func OperWrite(handle int32, data []byte, n int) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()

	if s.b != resAssigned || s.s == nil {
		return errors.New(ErrUnknownSocket)
	}

	if s.s.getStat() != Connected {
		return errors.New(ErrUnConnected)
	}

	s.s.push(&NetChunk{Handle: handle, Data: data}, n)

	return nil
}

// OperUDPWrite : 发送udp数据
func OperUDPWrite(handle int32, addr string, data []byte, n int) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()

	if s.b != resAssigned || s.s == nil {
		return errors.New(ErrUnknownSocket)
	}

	if s.s.getStat() != Connected {
		return errors.New(ErrUnConnected)
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	s.s.push(&NetChunk{Handle: handle, Data: data, Addr: udpAddr.IP, Port: uint16(udpAddr.Port)}, n)

	return nil
}

// OperKeep : 获取 socket 得Keep
func OperKeep(handle int32) (uint64, error) {
	s := defaultNServer.get(handle)
	if s == nil {
		return 0, errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return 0, errors.New(ErrUnknownSocket)
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return 0, errors.New(ErrClosedSocket)
	}
	return s.s.getKeepAive(), nil
}

// OperSetKeep : 设置socket keep
func OperSetKeep(handle int32, keep uint64) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return errors.New(ErrUnknownSocket)
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return errors.New(ErrClosedSocket)
	}

	s.s.setKeepAive(keep)
	s.s.getLastActivedTime()
	return nil
}

// OperLastActivedTime : 获取 socket 最后的活动事件
func OperLastActivedTime(handle int32) (uint64, error) {
	s := defaultNServer.get(handle)
	if s == nil {
		return 0, errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return 0, errors.New(ErrUnknownSocket)
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return 0, errors.New(ErrClosedSocket)
	}

	return s.s.getLastActivedTime(), nil
}

// OperOpen : 打开socket
func OperOpen(handle int32) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return errors.New(ErrUnknownSocket)
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return errors.New(ErrUnknownSocket)
	}

	if !s.s.setConnected() {
		return errors.New(ErrClosedSocket)
	}
	return nil
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
