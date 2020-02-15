package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/args"
)

//NetClient
type NetClient struct {
	_conn    io.ReadWriteCloser
	_queue   chan []byte
	_reader  *bufio.Reader
	_writer  *bufio.Writer
	_closed  chan bool
	_timeout time.Time
	_check   int
	_state   int
	_wg      sync.WaitGroup
}

func (slf *NetClient) Connect(addr string, timeout time.Duration) error {
	if slf._closed != nil {
		slf._closed = make(chan bool, 1)
	}
	c, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}

	slf._conn = c
	slf._reader = bufio.NewReaderSize(slf._conn, 8192)
	slf._writer = bufio.NewWriterSize(slf._conn, 8192)

	slf._wg.Add(1)
	go func() {
		for {
			_, err := slf.Parse()
			if err != nil {
				slf._closed <- true
				break
			}
		}
	}()

	go func() {
		defer func() {
			slf._state = 1
			slf._conn.Close()
			slf._wg.Done()
		}()

		for {
			select {
			case <-slf._closed:
				goto exit
			case msg := <-slf._queue:

				if err := slf.Write(msg); err != nil {
					fmt.Println("write error,", err)
					goto exit
				}
			}
		}
	exit:
	}()

	return err
}

func (slf *NetClient) Parse() (interface{}, error) {
	var header uint16
	if err := binary.Read(slf._reader, binary.BigEndian, &header); err != nil {
		return nil, err
	}

	length := int(header)
	buffer := make([]byte, length)
	offset := 0
	for offset < length {
		n, err := slf._reader.Read(buffer[offset:])
		if err != nil {
			return nil, err
		}

		offset += n

		if err != nil && offset < length {
			// if we read whole size of message, ignore error at this time.
			return nil, err
		}
	}

	return string(buffer), nil
}

func (slf *NetClient) Write(b []byte) error {
	length := len(b)
	seek := 0
	for {
		if slf._writer.Available() == 0 {
			if err := slf._writer.Flush(); err != nil {
				return err
			}
		}

		n, err := slf._writer.Write(b[seek:])
		if err != nil {
			return err
		}

		seek += n
		if seek >= length {
			break
		}
	}

	if slf._writer.Buffered() > 0 {
		if err := slf._writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (slf *NetClient) Push(data []byte) error {
	if slf._state != 0 {
		return errors.New("connection closed")
	}
	slf._queue <- data
	return nil
}

func (slf *NetClient) Close() {
	slf._state = 1
	if slf._conn != nil {
		slf._conn.Close()
	}

	if slf._closed != nil {
		slf._closed <- true
	}
}

func main() {
	args.Instance().Parse()
	addr := args.Instance().GetString("-p", "127.0.0.1:12000")
	maxConn := args.Instance().GetInt("-n", 1)
	deply := args.Instance().GetInt("-d", 300)
	checkNum := args.Instance().GetInt("-c", 1)
	timeOut := args.Instance().GetInt("-t", 1000)

	var connected int
	var sendCount int
	var failCount int
	var clients []*NetClient
	timeDeply := time.Duration(deply)
	for i := 0; i < maxConn; i++ {
		c := &NetClient{
			_queue:  make(chan []byte, 16),
			_closed: make(chan bool),
		}

		fmt.Print("开始连接:", addr)
		if err := c.Connect(addr, time.Duration(timeOut)*time.Millisecond); err != nil {
			c.Close()
			fmt.Println(" 连接失败 ", err)
			continue
		}

		fmt.Println(" 连接成功")

		c._timeout = time.Now()
		clients = append(clients, c)
		connected++
	}

	for {
		curtime := time.Now()
		//for i, cc := range clients {
		for i := 0; i < len(clients); {
			cc := clients[i]
			diff := curtime.Sub(cc._timeout)
			if diff.Milliseconds() >= int64(timeDeply) {
				if cc._check >= checkNum {
					clients = append(clients[0:i], clients[i+1:]...)
					cc.Close()
					continue
				}

				b := make([]byte, 2)
				binary.BigEndian.PutUint16(b, 4)
				b = append(b, []byte("abcd")...)
				if err := cc.Push(b); err != nil {
					failCount++
					clients = append(clients[0:i], clients[i+1:]...)
					cc.Close()
					continue
				}
				sendCount++
				cc._timeout = curtime
				cc._check++
			}
			i++
		}
		if len(clients) == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("总连接次数:", maxConn)
	fmt.Println("完成连接数:", connected)
	fmt.Println("发送成功次数:", sendCount)
	fmt.Println("发送失败次数:", failCount)
}
