package netmsgs

import "net"

//Accept netboxs connection accept message
type Accept struct {
	Sock int32
	Addr net.Addr
}

//Message netboxs connection recvice message
type Message struct {
	Sock interface{}
	Data interface{}
}

//Closed netboxs connection closed message
type Closed struct {
	Sock int32
}

//Error netboxs connecton error message
type Error struct {
	Sock interface{}
	Err  error
}
