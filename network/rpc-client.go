package network

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"net"
	"net/rpc"
	"time"
)

type rpcClientCodec struct {
	rwc      io.ReadWriteCloser
	dec      *gob.Decoder
	enc      *gob.Encoder
	encBuf   *bufio.Writer
	operator *actor.PID
}

func (c *rpcClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {
	if err = timeoutCoder(c.enc.Encode, r, "client write request"); err != nil {
		logger.Error(c.operator.ID, "client write request:%s", err.Error())
		return
	}
	if err = timeoutCoder(c.enc.Encode, body, "client write request body"); err != nil {
		logger.Error(c.operator.ID, "client write request body:%s", err.Error())
		return
	}
	return c.encBuf.Flush()
}

func (c *rpcClientCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.dec.Decode(r)
}

func (c *rpcClientCodec) ReadResponseBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *rpcClientCodec) Close() error {
	return c.rwc.Close()
}

type rpcClient struct {
	h        int32
	operator *actor.PID
	stat     int32
}

func (rpcc *rpcClient) call(operator *actor.PID, addr string, rpcmethod string, args interface{}, reply interface{}) error {
	conn, err := net.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return fmt.Errorf("connect error:%s", err.Error())
	}

	encBuf := bufio.NewWriter(conn)
	codec := &rpcClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf, operator}
	c := rpc.NewClientWithCodec(codec)
	err = c.Call(rpcmethod, args, reply)
	errc := c.Close()
	if err != nil && errc != nil {
		return fmt.Errorf("%s %s", err, errc)
	}

	if err != nil {
		return err
	}
	return errc
}
