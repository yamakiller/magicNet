package implement

import (
	"bytes"
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
	SetID(uint64)
	GetID() uint64
	GetAuth() uint64
	SetAuth(v uint64)
	GetSocket() int32
	SetSocket(sock int32)
	SetAddr(adr string)
	GetAddr() string
	GetRecvBuffer() *bytes.Buffer
	SetRecvBuffer(b *bytes.Buffer)
	GetKeyPair() interface{}
	BuildKeyPair()
	GetKeyPublic() string
	GetStat() *NetStat
	Shutdown()
	SetRef(v int)
	IncRef()
	DecRef() int
}

//NetClient Network client base class
type NetClient struct {
	wb   *bytes.Buffer
	addr string
	stat NetStat
	ref  int
}

//SetAddr Setting client address information
func (nc *NetClient) SetAddr(adr string) {
	nc.addr = adr
}

//GetAddr returns client address
func (nc *NetClient) GetAddr() string {
	return nc.addr
}

//GetRecvBuffer Return client read buffer
func (nc *NetClient) GetRecvBuffer() *bytes.Buffer {
	return nc.wb
}

//SetRecvBuffer Setting  the client read buffer
func (nc *NetClient) SetRecvBuffer(b *bytes.Buffer) {
	nc.wb = b
}

//SetRef Setting the number of citations
func (nc *NetClient) SetRef(v int) {
	nc.ref = v
}

//IncRef Add a reference
func (nc *NetClient) IncRef() {
	nc.ref++
}

//DecRef Reduce one reference
func (nc *NetClient) DecRef() int {
	nc.ref--
	return nc.ref
}

//GetStat returns status
func (nc *NetClient) GetStat() *NetStat {
	return &nc.stat
}

//Shutdown termination client
func (nc *NetClient) Shutdown() {

}
