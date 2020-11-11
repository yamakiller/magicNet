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
	OpenTo(int32) error
	SendTo(interface{}, []byte) error
	CloseTo(int32) error
	CloseToWait(int32) error
	GetConnect(int32) (interface{}, error)
	GetValues() []int32
}
