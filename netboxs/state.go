package netboxs

type state int32

const (
	stateInit state = iota
	stateConnecting
	stateConnected
	stateAccepted
	stateIdle
	stateDetached
	stateSend
	stateReceive
	stateShutdown
	stateClosed
)
