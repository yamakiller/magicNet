package network

import (
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicLibs/mutex"
	"github.com/yamakiller/magicNet/engine/actor"

	"google.golang.org/grpc"
)

// GRpcRegister : 注册 GRPC 方法
type GRpcRegister func(s *grpc.Server)

type grpcServer struct {
	h        int32
	s        *net.TCPListener
	rpc      *grpc.Server
	so       *slot
	operator *actor.PID
	stat     int32
	netWait  sync.WaitGroup

	register              GRpcRegister
	writeBufSize          int
	readBufSize           int
	connectionTimeout     int
	maxSendMessageSize    int
	maxReceiveMessageSize int
}

func (slf *grpcServer) listen(operator *actor.PID, addr string) error {
	if slf.writeBufSize == 0 {
		slf.writeBufSize = 32 * 1024
	}

	if slf.readBufSize == 0 {
		slf.readBufSize = 32 * 1024
	}

	if slf.connectionTimeout == 0 {
		slf.connectionTimeout = 120
	}

	if slf.maxSendMessageSize == 0 {
		slf.maxSendMessageSize = math.MaxInt32
	}

	if slf.maxReceiveMessageSize == 0 {
		slf.maxReceiveMessageSize = 1024 * 1024 * 4
	}

	tcpAddr, aderr := net.ResolveTCPAddr("tcp", addr)
	if aderr != nil {
		return aderr
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	slf.s = ln
	slf.rpc = grpc.NewServer(grpc.ReadBufferSize(slf.readBufSize),
		grpc.WriteBufferSize(slf.writeBufSize),
		grpc.MaxRecvMsgSize(slf.maxReceiveMessageSize),
		grpc.MaxSendMsgSize(slf.maxSendMessageSize),
		grpc.ConnectionTimeout(time.Duration(slf.connectionTimeout)*time.Second))
	slf.stat = Connecting

	if slf.register != nil {
		slf.register(slf.rpc)
	}

	slf.netWait.Add(1)
	go slf.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (slf *grpcServer) serve(ln net.Listener) {
	defer slf.netWait.Done()
	for {
		err := slf.rpc.Serve(slf.s)
		logger.Error(slf.operator.ID, "grpc server error %v", err)
		if err != grpc.ErrServerStopped {
			goto rpc_end
		} else if err == grpc.ErrClientConnClosing ||
			err == grpc.ErrClientConnTimeout {
			continue
		}

		break
	}

	slf.stat = Closing
	slf.rpc.Stop()
rpc_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	slf.so.l.Lock()

	closeHandle = slf.h
	closeOperator = slf.operator
	slf.s.Close()

	slf.so.s = nil
	slf.so.b = resIdle

	slf.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{Handle: closeHandle})
}

func (slf *grpcServer) connect(operator *actor.PID, addr string) error {
	return nil
}

func (slf *grpcServer) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (slf *grpcServer) push(data *NetChunk, n int) error {
	return nil
}

func (slf *grpcServer) recv() {

}

func (slf *grpcServer) write() {

}

func (slf *grpcServer) setKeepAive(keep uint64) {

}

func (slf *grpcServer) getKeepAive() uint64 {
	return 0
}

func (slf *grpcServer) getLastActivedTime() uint64 {
	return 0
}

func (slf *grpcServer) getStat() int32 {
	return slf.stat
}

func (slf *grpcServer) getProto() string {
	return protoRPC
}

func (slf *grpcServer) getType() int {
	return CListen
}

func (slf *grpcServer) setConnected() bool {
	return atomic.CompareAndSwapInt32(&slf.stat, Connecting, Connected)
}

func (slf *grpcServer) close(lck *mutex.ReSpinLock) {
	if slf.stat != Closing {
		slf.stat = Closing
		slf.rpc.Stop()
	}
}

func (slf *grpcServer) closewait() {
	slf.netWait.Wait()
}
