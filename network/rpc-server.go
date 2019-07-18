package network

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"net"
	"net/rpc"
	"sync"
	"sync/atomic"
	"time"
)

func timeoutCoder(f func(interface{}) error, e interface{}, msg string) error {
	echan := make(chan error, 1)
	go func() { echan <- f(e) }()
	select {
	case e := <-echan:
		return e
	case <-time.After(time.Minute):
		return fmt.Errorf("Timeout %s", msg)
	}
}

type rpcServerCodec struct {
	rwc      io.ReadWriteCloser
	dec      *gob.Decoder
	enc      *gob.Encoder
	encBuf   *bufio.Writer
	operator *actor.PID
	closed   bool
}

func (c *rpcServerCodec) ReadRequestHeader(r *rpc.Request) error {
	return timeoutCoder(c.dec.Decode, r, "server read request header")
}

func (c *rpcServerCodec) ReadRequestBody(body interface{}) error {
	return timeoutCoder(c.dec.Decode, body, "server read request body")
}

func (c *rpcServerCodec) WriteResponse(r *rpc.Response, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, "server write response"); err != nil {
		if c.encBuf.Flush() == nil {
			logger.Error(c.operator.ID, "rpc: codec error encoding response:%s", err.Error())
			c.Close()
		}
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, "server write response body"); err != nil {
		if c.encBuf.Flush() == nil {
			logger.Error(c.operator.ID, "rpc: codec error encoding body:%s", err.Error())
			c.Close()
		}
		return
	}
	return c.encBuf.Flush()
}

func (c *rpcServerCodec) Close() error {
	if c.closed {
		// Only call c.rwc.Close once; otherwise the semantics are undefined.
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}

type rpcServer struct {
	h        int32
	s        *net.TCPListener
	so       *slot
	operator *actor.PID
	stat     int32
	netWait  sync.WaitGroup
}

func (rpcs *rpcServer) listen(operator *actor.PID, addr string) error {
	tcpAddr, aderr := net.ResolveTCPAddr("tcp", addr)
	if aderr != nil {
		return aderr
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	rpcs.s = ln
	rpcs.stat = Connecting
	rpcs.netWait.Add(1)
	go rpcs.serve(ln)

	time.Sleep(time.Millisecond * 1)
	return nil
}

func (rpcs *rpcServer) serve(ln net.Listener) {
	defer rpcs.netWait.Done()
	for {
		conn, err := rpcs.s.AcceptTCP()
		if err != nil {
			break
		}

		rpcs.netWait.Add(1)
		go func(conn net.Conn, rpcs *rpcServer) {
			defer rpcs.netWait.Done()
			buf := bufio.NewWriter(conn)
			srv := &rpcServerCodec{
				rwc:      conn,
				dec:      gob.NewDecoder(conn),
				enc:      gob.NewEncoder(buf),
				encBuf:   buf,
				operator: rpcs.operator,
			}
			defer srv.Close()
			err = rpc.ServeRequest(srv)
			if err != nil {
				logger.Error(rpcs.operator.ID, "server rpc request", err.Error())
			}
		}(conn, rpcs)
	}

	rpcs.so.l.Lock()
	defer rpcs.so.l.Unlock()
	rpcs.close(nil)
	rpcs.so.s = nil
	rpcs.so.b = resIdle

}

func (rpcs *rpcServer) connect(operator *actor.PID, addr string) error {
	return nil
}

func (rpcs *rpcServer) udpConnect(operator *actor.PID, srcAddr string, dstAddr string) error {
	return nil
}

func (rpcs *rpcServer) push(data *NetChunk, n int) error {
	return nil
}

func (rpcs *rpcServer) recv() {

}

func (rpcs *rpcServer) write() {

}

func (rpcs *rpcServer) setKeepAive(keep uint64) {

}

func (rpcs *rpcServer) getKeepAive() uint64 {
	return 0
}

func (rpcs *rpcServer) getLastActivedTime() uint64 {
	return 0
}

func (rpcs *rpcServer) getStat() int32 {
	return rpcs.stat
}

func (rpcs *rpcServer) getProto() string {
	return protoRPC
}

func (rpcs *rpcServer) getType() int {
	return CListen
}

func (rpcs *rpcServer) setConnected() bool {
	return atomic.CompareAndSwapInt32(&rpcs.stat, Connecting, Connected)
}

func (rpcs *rpcServer) close(lck *util.ReSpinLock) {
	if rpcs.stat != Closing {
		rpcs.stat = Closing
		rpcs.s.Close()
	}
}

func (rpcs *rpcServer) closewait() {
	rpcs.netWait.Wait()
}
