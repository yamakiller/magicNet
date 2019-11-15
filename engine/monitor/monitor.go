package monitor

//IMonitor desc
//@interface IMonitor desc: Monitor model interface
//@method (IsShutdown() bool) Whether the system has been closed
//@method (Shutdown()) Shut down system
//@method (IncService()) Add a service
//@method (DecService()) Reduce a service
//@method (WaitService()) Waiting for all services to end
type IMonitor interface {
	IsShutdown() bool
	Shutdown()
	IncService()
	DecService()
	WaitService()
}

var systemMonitor IMonitor

//WithMonitor desc
//@method WithMonitor desc: Setting up the system monitor
//@param (IMonitor) monitor instance
func WithMonitor(m IMonitor) {
	systemMonitor = m
}

//IsShutdown desc
//@method IsShutdown desc: Whether the system has been terminated
//@return (bool)
func IsShutdown() bool {
	if systemMonitor == nil {
		return defMonitor.IsShutdown()
	}
	return systemMonitor.IsShutdown()
}

//Shutdown desc
//@method Shutdown desc: Termination system
func Shutdown() {
	if systemMonitor == nil {
		defMonitor.Shutdown()
		return
	}
	systemMonitor.Shutdown()
}

//IncService desc
//@method IncService desc: Add a service
func IncService() {
	if systemMonitor == nil {
		defMonitor.IncService()
		return
	}
	systemMonitor.IncService()
}

//DecService desc
//@method DecService desc: Reduce a service
func DecService() {
	if systemMonitor == nil {
		defMonitor.DecService()
		return
	}
	systemMonitor.DecService()
}

//WaitService desc
//@method WaitService desc: Waiting for all services to end
func WaitService() {
	if systemMonitor == nil {
		defMonitor.WaitService()
		return
	}
	systemMonitor.WaitService()
}
