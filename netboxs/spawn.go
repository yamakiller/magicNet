package netboxs

import (
	"errors"

	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/boxs"
)

var (
	//ModeTCPListener  tcp listener
	ModeTCPListener = "tcp listener"
	//ModeWSSListener  websocket listener
	ModeWSSListener = "wss listener"
	//ModeKCPListener  kcp(udp) listener
	ModeKCPListener = "kcp listener"
)

//Spawn create an network box
func Spawn(tag string, pool Pool) (actors.Actor, error) {
	switch tag {
	case ModeTCPListener:
		tcp := &TCPBox{
			Box: *boxs.SpawnBox(nil),
		}
		tcp.WithPool(pool)
		return tcp, nil
	case ModeWSSListener:
		wss := &WSSBox{
			Box: *boxs.SpawnBox(nil),
		}
		wss.WithPool(pool)
		return wss, nil
	default:
		return nil, errors.New("undefined")
	}
}
