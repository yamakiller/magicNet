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
func (c *TCPConnection) Name() string {
	return "TCP/Connection"
}

//Socket Returns the TCP connection socket
func (c *TCPConnection) Socket() int32 {
	return c.s
}

//Connection TCP connection remote server returns error message if connection fails
func (c *TCPConnection) Connection(context actor.Context, addr string, outChanSize int) error {
	sock, err := network.OperTCPConnect(context.Self(), addr, outChanSize)
	if err != nil {
		return err
	}
	c.s = sock
	return nil
}

//Write Send data to a remote server
func (c *TCPConnection) Write(wrap []byte, length int) error {
	return network.OperWrite(c.s, wrap, length)
}

//Close Close the connection
func (c *TCPConnection) Close() {
	if c.s == 0 {
		return
	}

	network.OperClose(c.s)
}
