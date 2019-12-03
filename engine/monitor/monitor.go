package monitor

//IMonitor desc
//@Interface IMonitor desc: Monitor model interface
//@Method (IsShutdown() bool) Whether the system has been closed
//@Method (Shutdown()) Shut down system
//@Method (IncService()) Add a service
//@Method (DecService()) Reduce a service
//@Method (WaitService()) Waiting for all services to end
type IMonitor interface {
	IsShutdown() bool
	Shutdown()
	IncService()
	DecService()
	WaitService()
}

var systemMonitor IMonitor

//WithMonitor desc
//@Method WithMonitor desc: Setting up the system monitor
//@Param (IMonitor) monitor instance
func WithMonitor(m IMonitor) {
	systemMonitor = m
}

//IsShutdown desc
//@Method IsShutdown desc: Whether the system has been terminated
//@Return (bool)
func IsShutdown() bool {
	if systemMonitor == nil {
		return defMonitor.IsShutdown()
	}
	return systemMonitor.IsShutdown()
}

//Shutdown desc
//@Method Shutdown desc: Termination system
func Shutdown() {
	if systemMonitor == nil {
		defMonitor.Shutdown()
		return
	}
	systemMonitor.Shutdown()
}

//IncService desc
//@Method IncService desc: Add a service
func IncService() {
	if systemMonitor == nil {
		defMonitor.IncService()
		return
	}
	systemMonitor.IncService()
}

//DecService desc
//@Method DecService desc: Reduce a service
func DecService() {
	if systemMonitor == nil {
		defMonitor.DecService()
		return
	}
	systemMonitor.DecService()
}

//WaitService desc
//@Method WaitService desc: Waiting for all services to end
func WaitService() {
	if systemMonitor == nil {
		defMonitor.WaitService()
		return
	}
	systemMonitor.WaitService()
}
