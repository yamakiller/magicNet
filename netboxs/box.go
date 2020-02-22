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
	OpenTo(socket int32) error
	SendTo(socket int32, data []byte) error
	CloseTo(socket int32) error
	CloseToWait(socket int32) error
	GetValues() []int32
}
