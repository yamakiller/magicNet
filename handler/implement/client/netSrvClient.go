package client

import (
	"github.com/yamakiller/magicNet/handler/net"
)

type NetSrvClient struct {
	Receive net.INetBuffer

	_addr            string
	_sock            int32
	_receiveBytes    int64
	_receiveLastTime int64
	_sendoutBytes    int64
	_sendoutLastTime int64
	_onine           int64
	_ref             int
}

func (slf *NetSrvClient) WithAddr(addr string) {
	slf._addr = addr
}

//WithSocket doc
//@Summary Setting the client socket
//@Method WithSocket
//@Param (int32) socket id
func (slf *NetSrvClient) WithSocket(sock int32) {
	slf._sock = sock
}

//GetSocket doc
//@Summary Returns the client socket
//@Method GetSocket
//@Return (int32) socket id
func (slf *NetSrvClient) GetSocket() int32 {
	return slf._sock
}

//GetStatistics doc
//@Summary Returns the client Statistics informat
//@Method GetStatistics
//@Return int64 receive bytes count
//@Return int64 receive last time
//@Return int64 sendto bytes count
//@Return int64 sendto last time
//@Return int64 time online
func (slf *NetSrvClient) GetStatistics() (recvBytes int64,
	recvLastTime int64,
	sendToBytes int64,
	sendToLastTime int64,
	online int64) {
	return slf._receiveBytes, slf._receiveLastTime, slf._sendoutBytes, slf._sendoutLastTime, slf._onine
}

//GetBufferCap doc
//@Summary Returns Recvice buffer cap
//@Method GetBufferCap
//@Return int
func (slf *NetSrvClient) GetBufferCap() int {
	return slf.Receive.Cap()
}

//GetBufferLen doc
//@Summary Returns Recvice buffer data length
//@Method GetBufferLen
//@Return int
func (slf *NetSrvClient) GetBufferLen() int {
	return slf.Receive.Len()
}

//ClearBuffer doc
//@Summary Clear Recvice buffer data
//@Method ClearBuffer
func (slf *NetSrvClient) ClearBuffer() {
	slf.Receive.Clear()
}

//WriteBuffer doc
//@Summary Write Recvice buffer data
//@Method WriteBuffer
//@Return int
//@Return error
func (slf *NetSrvClient) WriteBuffer(b []byte) (int, error) {
	return slf.Receive.Write(b)
}

//UpdateReceive doc
//@Summary update receive informat
//@Method UpdateReceive
//@Param int64 now time
//@Param int64 receive data length
func (slf *NetSrvClient) UpdateReceive(tm int64, bytes int64) {
	slf._receiveLastTime = tm
	slf._receiveBytes += bytes
}

//UpdateSendto doc
//@Summary update send to informat
//@Method UpdateSendto
//@Param int64 now time
//@Param int64 send to data length
func (slf *NetSrvClient) UpdateSendto(tm int64, bytes int64) {
	slf._sendoutLastTime = tm
	slf._sendoutBytes += bytes
}

//UpdateOnline doc
//@Summary update online
//@Method UpdateOnline
//@Param  int64 now time
func (slf *NetSrvClient) UpdateOnline(tm int64) {
	slf._onine = tm
}

//SetRef doc
//@Summary Setting the number of citations
//@Method SetRef
//@Param  int value
func (slf *NetSrvClient) SetRef(v int) {
	slf._ref = v
}

//IncRef doc
//@Summary Add a reference
//@Method IncRef
func (slf *NetSrvClient) IncRef() {
	slf._ref++
}

//DecRef doc
//@Summary Reduce one reference
//@Method DecRef
//@Return int Subtracted by 1
func (slf *NetSrvClient) DecRef() int {
	slf._ref--
	return slf._ref
}
