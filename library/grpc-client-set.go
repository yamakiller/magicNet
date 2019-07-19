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
func (set *GRPCClientSet) Init(target string, opts ...grpc.DialOption) error {
	set.conns = make(chan *grpcIdleConn, set.MaxConn)
	set.factory = func() (*grpc.ClientConn, error) {
		ctx, cancel := context.WithTimeout(context.Background(), set.ConnTimeout)
		defer cancel()

		return grpc.DialContext(ctx, target)
	}
	set.close = func(c *grpc.ClientConn) error { return c.Close() }

	for i := 0; i < set.MinConn; i++ {
		conn, err := set.factory()
		if err != nil {
			set.Close()
			return nil
		}

		set.conns <- &grpcIdleConn{conn: conn, t: time.Now()}
	}

	return nil
}

//Close : close set
func (set *GRPCClientSet) Close() {
	set.mx.Lock()
	conns := set.conns
	set.conns = nil
	set.factory = nil
	closeFun := set.close
	set.close = nil
	set.mx.Unlock()

	if conns == nil {
		return
	}

	close(conns) //？修改退出 未必可以全部删除
	for wrapConn := range conns {
		closeFun(wrapConn.conn)
	}
}

// Invoke : 调用方法
func (set *GRPCClientSet) Invoke(method string, args, reply interface{}) error { //优化参数设置
	conn, err := set.getConn()
	if err != nil {
		return err
	}
	defer set.putConn(conn)

	ctx, cancel := context.WithTimeout(context.Background(), set.ConnTimeout)
	defer cancel()

	return conn.Invoke(ctx, method, args, reply)
}

func (set *GRPCClientSet) getConn() (*grpc.ClientConn, error) {
	set.mx.Lock()
	conns := set.conns
	set.mx.Unlock()

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
			if timeout := set.IdleTimeout; timeout > 0 {
				if wrapConn.t.Add(timeout).Before(time.Now()) {
					//丢弃并关闭该链接
					set.close(wrapConn.conn)
					continue
				}
			}
			return wrapConn.conn, nil
		default:
			conn, err := set.factory()
			if err != nil {
				return nil, err
			}

			return conn, nil
		}
	}
}

func (set *GRPCClientSet) putConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return errGRPCSetRejected
	}

	set.mx.Lock()
	defer set.mx.Unlock()

	if set.conns == nil {
		return set.close(conn)
	}

	select {
	case set.conns <- &grpcIdleConn{conn: conn, t: time.Now()}:
		return nil
	default:
		//连接池已满，直接关闭该链接
		return set.close(conn)
	}
}

// IdleCount : 空闲连接数
func (set *GRPCClientSet) IdleCount() int {
	set.mx.Lock()
	conns := set.conns
	set.mx.Unlock()
	return len(conns)
}
