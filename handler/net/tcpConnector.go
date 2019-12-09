package net

import (
	"github.com/yamakiller/magicLibs/net"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
)

//TCPConnection TCP connection object
type TCPConnection struct {
	_s int32
}

//GetSocket Returns the TCP connection socket
func (slf *TCPConnection) GetSocket() int32 {
	return slf._s
}

//Connection TCP connection remote server returns error message if connection fails
func (slf *TCPConnection) Connection(context actor.Context, addr string, outChanSize int) error {
	sock, err := network.OperTCPConnect(context.Self(), addr, outChanSize)
	if err != nil {
		return err
	}
	slf._s = sock
	return nil
}

//Write Send data to a remote server
func (slf *TCPConnection) Write(wrap []byte, length int) error {
	return network.OperWrite(slf._s, wrap, length)
}

//Close Close the connection
func (slf *TCPConnection) Close() {
	if net.InvalidSocket(slf._s) {
		return
	}

	network.OperClose(slf._s)
	slf._s = net.INVALIDSOCKET
}

//ToString doc
//@Summary to string
//@Method ToString
//@Return  string
func (slf *TCPConnection) ToString() string {
	return "TCP/Conn"
}
