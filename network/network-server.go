package network

import (
	"errors"
	"net"
	"sync/atomic"

	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"
)

const (
	maxSocket = 65535
)

const (
	resIdle     = 0
	resOccupy   = 1
	resAssigned = 2
)

var (
	// ErrSocketResources : Insufficient socket resources
	ErrSocketResources = errors.New("Insufficient socket resources")
	// ErrUnknownSocket :
	ErrUnknownSocket = errors.New("unknown socket")
	// ErrClosedSocket :
	ErrClosedSocket = errors.New("socket is closed")
	// ErrUnConnected :
	ErrUnConnected = errors.New("socket is not connected")
)

// OperWSListen : Open a websocket listening service and return the handle corresponding to the socket
func OperWSListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	s.l.Lock()
	s.s = &wsServer{}
	wss, _ := s.s.(*wsServer)
	wss._h = h
	wss._operator = operator
	wss._outChanMax = outChanSize

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

// OperWSConnect : Open an external connection to a websocket and return the handle corresponding to the socket
func OperWSConnect(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	client := &wsClient{}
	client._h = h
	client._o = operator
	client._out = make(chan *NetChunk, outChanSize)
	client._so = s
	err := client.connect(operator, addr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}

	s.l.Lock()
	s.s = client
	client._w.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperTCPListen : Open a socket TCP listening service and return the handle of the socket
func OperTCPListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	s.l.Lock()
	s.s = &tcpServer{}
	tps, _ := s.s.(*tcpServer)
	tps._h = h
	tps._operator = operator
	tps._outChanMax = outChanSize

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

// OperTCPConnect : Open a TCP connection to a socket and return the handle of the socket
func OperTCPConnect(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	client := &tcpClient{}
	client._h = h
	client._o = operator
	client._out = make(chan *NetChunk, outChanSize)
	client._quit = make(chan int)
	client._so = s
	err := client.connect(operator, addr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}

	s.l.Lock()
	s.s = client
	client._w.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperUDPListen : Open a socket UDP listening service and return the handle of the socket
func OperUDPListen(operator *actor.PID, addr string, outChanSize int) (int32, error) {
	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	s.l.Lock()
	s.s = &udpSocket{}
	ups, _ := s.s.(*udpSocket)
	ups._h = h
	ups._so = s
	ups._operator = operator
	ups._out = make(chan *NetChunk, outChanSize)
	ups._quit = make(chan int)

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

// OperUDPConnect : Open the UDP client of a socket and return the handle of the socket
func OperUDPConnect(operator *actor.PID, srcAddr string, dstAddr string, outChanSize int) (int32, error) {

	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
	}

	client := &udpSocket{}
	client._h = h
	client._operator = operator
	client._out = make(chan *NetChunk, outChanSize)
	client._quit = make(chan int)
	client._so = s

	err := client.udpConnect(operator, srcAddr, dstAddr)
	if err != nil {
		atomic.StoreInt32(&s.b, resIdle)
		return -1, err
	}
	s.l.Lock()
	s.s = client
	client._netWait.Add(2)
	go client.recv()
	go client.write()
	atomic.StoreInt32(&s.b, resAssigned)
	s.l.Unlock()

	return h, nil
}

// OperRPCListen Start an RPC service and return to handle. This RPC is based on gRPC
func OperRPCListen(operator *actor.PID,
	addr string,
	reg GRpcRegister, /*(Registration Protocol Call Function)*/
	writeBufSize int, /*(Optional default: 32 * 1024)*/
	readBufSize int, /*(Optional default: 32 * 1024)*/
	connectionTimeout int, /*(Optional default: 120 sec)*/
	maxSendMessageSize int, /*(Optional default: math.MaxInt32)*/
	maxReceiveMessageSize int /*(Optional default: 1024 * 1024 * 4)*/) (int32, error) {

	h, s := defaultNServer.grap()
	if h == -1 || s == nil {
		return h, ErrSocketResources
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

//OperWrite : Send data to the socket corresponding to handle
func OperWrite(handle int32, data []byte, n int) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()

	if s.b != resAssigned || s.s == nil {
		return ErrUnknownSocket
	}

	if s.s.getStat() != Connected {
		return ErrUnConnected
	}

	s.s.push(&NetChunk{Handle: handle, Data: data}, n)

	return nil
}

// OperUDPWrite : Send UDP data to the socket corresponding to handle
func OperUDPWrite(handle int32, addr string, data []byte, n int) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()

	if s.b != resAssigned || s.s == nil {
		return ErrUnknownSocket
	}

	if s.s.getStat() != Connected {
		return ErrUnConnected
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	s.s.push(&NetChunk{Handle: handle, Data: data, Addr: udpAddr.IP, Port: uint16(udpAddr.Port)}, n)

	return nil
}

// OperKeep : Returns Socket Keep information corresponding to handle
func OperKeep(handle int32) (uint64, error) {
	s := defaultNServer.get(handle)
	if s == nil {
		return 0, ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return 0, ErrUnknownSocket
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return 0, ErrClosedSocket
	}
	return s.s.getKeepAive(), nil
}

// OperSetKeep : Setting keep of handle corresponding socket
func OperSetKeep(handle int32, keep uint64) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return ErrUnknownSocket
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return ErrClosedSocket
	}

	s.s.setKeepAive(keep)
	s.s.getLastActivedTime()
	return nil
}

// OperLastActivedTime : Returns the last active time of the handle corresponding socket
func OperLastActivedTime(handle int32) (uint64, error) {
	s := defaultNServer.get(handle)
	if s == nil {
		return 0, ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return 0, ErrUnknownSocket
	}

	if s.s.getStat() != Connecting && s.s.getStat() != Connected {
		return 0, ErrClosedSocket
	}

	return s.s.getLastActivedTime(), nil
}

// OperOpen : Modify the connected state of the socket corresponding to handle
func OperOpen(handle int32) error {
	s := defaultNServer.get(handle)
	if s == nil {
		return ErrUnknownSocket
	}

	s.l.Lock()
	defer s.l.Unlock()
	if s.b != resAssigned || s.s == nil {
		return ErrUnknownSocket
	}

	if !s.s.setConnected() {
		return ErrClosedSocket
	}
	return nil
}

// OperClose : Close the socket corresponding to handle
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
	l mutex.SpinLock
}

// NetServer : Network Server Manager
type NetServer struct {
	ss []slot
	fi int32
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
