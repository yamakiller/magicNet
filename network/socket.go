package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/util"
)

const (
	//ProtoTCP : TCP/IP　协议
	ProtoTCP = "tcp"
	//ProtoUDP : UDP　协议
	ProtoUDP = "udp"
	//ProtoWeb : WebSocket
	ProtoWeb = "web"
)

const (
	// CListen : 监听
	CListen = iota
	// CConnect : 连接
	CConnect
)

const (
	// Idle : 空闲
	Idle = iota
	// Connecting : 连接中
	Connecting
	// Connected : 已经连接
	Connected
	// Closing  : 关闭中
	Closing
	// Closed    : 已经不安比
	Closed
)

// NetChunk : 网络数据消息
type NetChunk struct {
	Data []byte
}

// NetClose : Socket 关闭消息
type NetClose struct {
	handle int32
}

// NetInfo  : Socket 通用状态数据信息
type NetInfo struct {
	WriteBytes    uint64
	WriteLastTime uint64
	ReadBytes     uint64
	ReadLastTime  uint64
}

// ISocket : 套接字接口
type ISocket interface {
	listen(operator *actor.PID, addr string) error
	connect(operator *actor.PID, addr string) error
	setKeepAive(keep uint64)
	getStat() int
	setStat(stat int)
	close(lck *util.ReSpinLock)
	closewait()
}
