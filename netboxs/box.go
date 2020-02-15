package netboxs

//NetBox network connection or listener server
type NetBox interface {
	Shutdown()
}

//NetLBox network listener server
type NetLBox interface {
	NetBox
	WithPool(Pool)
	WithMax(int32)
	ListenAndServe(addr string) error
	OpenConn(socket int32) error
	SendConn(socket int32, data []byte) error
	CloseConn(socket int32) error
	CloseConnWait(socket int32) error
	GetValues() []int32
}
