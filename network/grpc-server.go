package network

import (
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/util"

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

func (rpcs *grpcServer) listen(operator *actor.PID, addr string) error {
	if rpcs.writeBufSize == 0 {
		rpcs.writeBufSize = 32 * 1024
	}

	if rpcs.readBufSize == 0 {
		rpcs.readBufSize = 32 * 1024
	}

	if rpcs.connectionTimeout == 0 {
		rpcs.connectionTimeout = 120
	}

	if rpcs.maxSendMessageSize == 0 {
		rpcs.maxSendMessageSize = math.MaxInt32
	}

	if rpcs.maxReceiveMessageSize == 0 {
		rpcs.maxReceiveMessageSize = 1024 * 1024 * 4
	}

	tcpAddr, aderr := net.ResolveTCPAddr("tcp", addr)
	if aderr != nil {
		return aderr
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	rpcs.s = ln
	rpcs.rpc = grpc.NewServer(grpc.ReadBufferSize(rpcs.readBufSize),
		grpc.WriteBufferSize(rpcs.writeBufSize),
		grpc.MaxRecvMsgSize(rpcs.maxReceiveMessageSize),
		grpc.MaxSendMsgSize(rpcs.maxSendMessageSize),
		grpc.ConnectionTimeout(time.Duration(rpcs.connectionTimeout)*time.Second))
	rpcs.stat = Connecting

	if rpcs.register != nil {
		rpcs.register(rpcs.rpc)
	}

	rpcs.netWait.Add(1)
	go rpcs.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (rpcs *grpcServer) serve(ln net.Listener) {
	defer rpcs.netWait.Done()
	for {
		err := rpcs.rpc.Serve(rpcs.s)
		logger.Error(rpcs.operator.ID, "grpc server error %v", err)
		if err != grpc.ErrServerStopped {
			goto rpc_end
		} else if err == grpc.ErrClientConnClosing ||
			err == grpc.ErrClientConnTimeout {
			continue
		}

		break
	}

	rpcs.stat = Closing
	rpcs.rpc.Stop()
rpc_end:
	var (
		closeHandle   int32
		closeOperator *actor.PID
	)

	rpcs.so.l.Lock()

	closeHandle = rpcs.h
	closeOperator = rpcs.operator
	rpcs.s.Close()

	rpcs.so.s = nil
	rpcs.so.b = resIdle

	rpcs.so.l.Unlock()

	actor.DefaultSchedulerContext.Send(closeOperator, NetClose{Handle: closeHandle})
}

func (rpcs *grpcServer) connect(operator *actor.PID, addr string) error {
	return nil
}

func (rpcs *grpcServer) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (rpcs *grpcServer) push(data *NetChunk, n int) error {
	return nil
}

func (rpcs *grpcServer) recv() {

}

func (rpcs *grpcServer) write() {

}

func (rpcs *grpcServer) setKeepAive(keep uint64) {

}

func (rpcs *grpcServer) getKeepAive() uint64 {
	return 0
}

func (rpcs *grpcServer) getLastActivedTime() uint64 {
	return 0
}

func (rpcs *grpcServer) getStat() int32 {
	return rpcs.stat
}

func (rpcs *grpcServer) getProto() string {
	return protoRPC
}

func (rpcs *grpcServer) getType() int {
	return CListen
}

func (rpcs *grpcServer) setConnected() bool {
	return atomic.CompareAndSwapInt32(&rpcs.stat, Connecting, Connected)
}

func (rpcs *grpcServer) close(lck *util.ReSpinLock) {
	if rpcs.stat != Closing {
		rpcs.stat = Closing
		rpcs.rpc.Stop()
	}
}

func (rpcs *grpcServer) closewait() {
	rpcs.netWait.Wait()
}
