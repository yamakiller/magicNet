package net

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPConnection TCP connection object
type TCPConnection struct {
	s int32
}

//Name Object name
func (slf *TCPConnection) Name() string {
	return "TCP/Connection"
}

//Socket Returns the TCP connection socket
func (slf *TCPConnection) Socket() int32 {
	return slf.s
}

//Connection TCP connection remote server returns error message if connection fails
func (slf *TCPConnection) Connection(context actor.Context, addr string, outChanSize int) error {
	sock, err := network.OperTCPConnect(context.Self(), addr, outChanSize)
	if err != nil {
		return err
	}
	slf.s = sock
	return nil
}

//Write Send data to a remote server
func (slf *TCPConnection) Write(wrap []byte, length int) error {
	return network.OperWrite(slf.s, wrap, length)
}

//Close Close the connection
func (slf *TCPConnection) Close() {
	if slf.s == 0 {
		return
	}

	network.OperClose(slf.s)
}
