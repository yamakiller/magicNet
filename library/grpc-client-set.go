package library

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type grpcIdleConn struct {
	conn *grpc.ClientConn
	t    time.Time
}

var (
	errGRPCSetClosed   = errors.New("pool is closed")
	errGRPCSetInvalid  = errors.New("invalid config")
	errGRPCSetRejected = errors.New("connection is nil. rejecting")
	errGRPCSetTargets  = errors.New("targets server is empty")
)

// GRPCClientSet : GRPC 客户端集
type GRPCClientSet struct {
	mx          sync.Mutex
	IdleTimeout time.Duration
	MinConn     int
	MaxConn     int
	ConnTimeout time.Duration
	conns       chan *grpcIdleConn
	factory     func() (*grpc.ClientConn, error)
	close       func(*grpc.ClientConn) error
}

//Init : 初始化客户端集
func (slf *GRPCClientSet) Init(target string, opts ...grpc.DialOption) error {
	slf.conns = make(chan *grpcIdleConn, slf.MaxConn)
	slf.factory = func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), slf.ConnTimeout)
		defer cancel()

		return grpc.DialContext(ctx, target)
	}
	slf.close = func(c *grpc.ClientConn) error { return c.Close() }

	for i := 0; i < slf.MinConn; i++ {
		conn, err := slf.factory()
		if err != nil {
			slf.Close()
			return nil
		}

		slf.conns <- &grpcIdleConn{conn: conn, t: time.Now()}
	}

	return nil
}

//Close : close set
func (slf *GRPCClientSet) Close() {
	slf.mx.Lock()
	conns := slf.conns
	slf.conns = nil
	slf.factory = nil
	closeFun := slf.close
	slf.close = nil
	slf.mx.Unlock()

	if conns == nil {
		return
	}

	close(conns) //？修改退出 未必可以全部删除
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

// Invoke : 调用方法
func (slf *GRPCClientSet) Invoke(method string, args, reply interface{}) error { //优化参数设置
	conn, err := slf.getConn()
	if err != nil {
		return err
	}
	defer slf.putConn(conn)

	ctx, cancel := context.WithTimeout(context.Background(), slf.ConnTimeout)
	defer cancel()

	return conn.Invoke(ctx, method, args, reply)
}

func (slf *GRPCClientSet) getConn() (*grpc.ClientConn, error) {
	slf.mx.Lock()
	conns := slf.conns
	slf.mx.Unlock()

	if conns == nil {
		return nil, errGRPCSetClosed
	}
	for {
		select {
		case wrapConn := <-conns:
			if wrapConn == nil {
				return nil, errGRPCSetClosed
			}
			//判断是否超时，超时则丢弃
			if timeout := slf.IdleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					slf.close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			conn, err := slf.factory()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}
	}
}

func (slf *GRPCClientSet) putConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return errGRPCSetRejected
	}

	slf.mx.Lock()
	defer slf.mx.Unlock()

	if slf.conns == nil {
		return slf.close(conn)
	}

	select {
	case slf.conns <- &grpcIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return slf.close(conn)
	}
}

// IdleCount : 空闲连接数
func (slf *GRPCClientSet) IdleCount() int {
	slf.mx.Lock()
	conns := slf.conns
	slf.mx.Unlock()
	return len(conns)
}
