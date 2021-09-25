package boxs

import "github.com/yamakiller/magicLibs/actors"

//Container implement container interface
type Container interface {
	GetPID() *actors.PID
	Register(interface{}, Method)
	StartedWait()
	Shutdown()
	ShutdownWait()
}
