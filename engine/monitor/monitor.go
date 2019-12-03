package monitor

//IMonitor doc
//@Interface IMonitor @Summary Monitor model interface
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

//WithMonitor doc
//@Method WithMonitor @Summary Setting up the system monitor
//@Param (IMonitor) monitor instance
func WithMonitor(m IMonitor) {
	systemMonitor = m
}

//IsShutdown doc
//@Method IsShutdown @Summary Whether the system has been terminated
//@Return (bool)
func IsShutdown() bool {
	if systemMonitor == nil {
		return defMonitor.IsShutdown()
	}
	return systemMonitor.IsShutdown()
}

//Shutdown doc
//@Method Shutdown @Summary Termination system
func Shutdown() {
	if systemMonitor == nil {
		defMonitor.Shutdown()
		return
	}
	systemMonitor.Shutdown()
}

//IncService doc
//@Method IncService @Summary Add a service
func IncService() {
	if systemMonitor == nil {
		defMonitor.IncService()
		return
	}
	systemMonitor.IncService()
}

//DecService doc
//@Method DecService @Summary Reduce a service
func DecService() {
	if systemMonitor == nil {
		defMonitor.DecService()
		return
	}
	systemMonitor.DecService()
}

//WaitService doc
//@Method WaitService @Summary Waiting for all services to end
func WaitService() {
	if systemMonitor == nil {
		defMonitor.WaitService()
		return
	}
	systemMonitor.WaitService()
}
