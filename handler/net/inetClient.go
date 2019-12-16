package net

import (
	"errors"
)

var (
	//ErrAnalysisSuccess Analyze data successfully
	ErrAnalysisSuccess = errors.New("Packet decomposition continues correctly")
	//ErrAnalysisProceed Not broken down to full data, please continue
	ErrAnalysisProceed = errors.New("No complete data package")
)

//INetClient Network client interface
type INetClient interface {
	WithID(uint64)
	GetID() uint64
	WithSocket(sock int32)
	WithAddr(addr string)
	GetSocket() int32
	GetAddr() string
	INetReceiveBuffer
	/*GetBufferCap() int
	GetBufferLen() int
	GetBufferBytes() []byte*/
	GetStatistics() (recvBytes int64, recvLastTime int64, sendToBytes int64, sendToLastTime int64, online int64)
	/*ClearBuffer()
	TrunBuffer(n int)
	WriteBuffer(b []byte) (int, error)
	ReadBuffer(n int) []byte*/
	SendTo(b []byte) error

	UpdateReceive(int64, int64)
	UpdateSendto(int64, int64)
	UpdateOnline(int64)
	Shutdown()
	SetRef(v int)
	IncRef()
	DecRef() int
}
