package netboxs

import "crypto/tls"

//NetBox network connection or listener server
type INetBox interface {
	Shutdown()
}

//NetLBox network listener server
type NetBox interface {
	INetBox
	WithPool(Pool)
	WithMax(int32)
	ListenAndServe(addr string) error
	ListenAndServeTls(addr string, ptls *tls.Config) error
	OpenTo(interface{}) error
	SendTo(interface{}, interface{}) error
	CloseTo(int32) error
	CloseToWait(int32) error
	GetConnect(int32) (interface{}, error)
	GetValues() []int32
}
