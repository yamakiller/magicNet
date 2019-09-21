package implement

import (
	"bytes"
	"errors"

	"github.com/yamakiller/magicNet/util"
)

var (
	//ErrAnalysisSuccess Analyze data successfully
	ErrAnalysisSuccess = errors.New("Packet decomposition continues correctly")
	//ErrAnalysisProceed Not broken down to full data, please continue
	ErrAnalysisProceed = errors.New("No complete data package")
)

//NetDataStat Network data status information
type NetDataStat struct {
	lastTime  uint64
	lastBytes uint64
}

//Update Update Network data status information
func (ndst *NetDataStat) Update(tts uint64, bytes uint64) {
	ndst.lastTime = tts
	ndst.lastBytes += bytes
}

//GetTime returns last time
func (ndst *NetDataStat) GetTime() uint64 {
	return ndst.lastTime
}

//GetBytes returns count bytes
func (ndst *NetDataStat) GetBytes() uint64 {
	return ndst.lastBytes
}

//NetStat network status
type NetStat struct {
	online uint64
	read   NetDataStat
	write  NetDataStat
}

//UpdateRead Update read data status
func (nst *NetStat) UpdateRead(tts uint64, bytes uint64) {
	nst.read.Update(tts, bytes)
}

//UpdateWrite Update write data status
func (nst *NetStat) UpdateWrite(tts uint64, bytes uint64) {
	nst.write.Update(tts, bytes)
}

//UpdateOnline Update online time information
func (nst *NetStat) UpdateOnline(tts uint64) {
	nst.online = tts
}

//GetRead returns read status object
func (nst *NetStat) GetRead() NetDataStat {
	return nst.read
}

//GetWrite returns write status object
func (nst *NetStat) GetWrite() NetDataStat {
	return nst.write
}

//GetOnline returns online time last
func (nst *NetStat) GetOnline() uint64 {
	return nst.online
}

//INetClient Network client interface
type INetClient interface {
	SetID(h *util.NetHandle)
	GetID() *util.NetHandle
	GetAuth() uint64
	SetAuth(v uint64)
	GetSocket() int32
	SetSocket(sock int32)
	SetAddr(adr string)
	GetAddr() string
	GetRecvBuffer() *bytes.Buffer
	SetRecvBuffer(b *bytes.Buffer)
	GetStat() *NetStat
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
