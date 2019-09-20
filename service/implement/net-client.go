package implement

import (
	"errors"
	"bytes"
)

var (
	//AnalysisSuccess
	AnalysisSuccess = errors.New("Packet decomposition continues correctly")
	//AnalysisProceed
	AnalysisProceed      = errors.Net("No complete data package")
)

type NetDataStat struct {
	lastTime  uint64
	lastBytes uint64
}

func (ndst *NetDataStat) Update(tts uint64, bytes uint64) {
	ndst.lastTime = tts
	ndst.lastBytes += bytes
}

func (ndst *NetDataStat) GetTime() uint64 {
	return ndst.lastTime
}

func (ndst *NetDataStat) GetBytes() uint64 {
	return ndst.lastBytes
}

type NetStat struct {
	online uint64
	read   NetDataStat
	write  NetDataStat
}

func (nst *NetStat) UpdateRead(tts uint64, bytes uint64) {
	nst.read.Update(tts, bytes)
}

func (nst *NetStat) UpdateWrite(tts uint64, bytes uint64) {
	nst.write.Update(tts, bytes)
}

func (nst *NetStat) UpdateOnline(tts uint64) {
	nst.online = tts
}

func (nst *NetStat) GetRead() NetDataStat {
	return nst.read
}

func (nst *NetStat) GetWrite() NetDataStat {
	return nst.write
}

func (nst *NetStat) GetOnline() uint64 {
	return nst.online
}

type INetClient interface {
	SetID(h *util.NetHandle)
	GetAuth() uint64
	SetAuth(v uint64)
	GetSocket() int32
	SetSocket(sock int32)
	SetAddr(ipaddr string)
	SetPort(port int)
	GetRecvBuffer() *bytes.Buffer
	SetRecvBuffer(b *bytes.Buffer)
	GetStat() *NetStat
	IncRef()
	DecRef() int
}
