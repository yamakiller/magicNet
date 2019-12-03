package monitor

//DefaultMonitor desc
//@Struct defaultMonitor desc: Default monitor
type defaultMonitor struct {
	isShutdown bool
}

//IsShutdown Whether the system has been closed
func (dmt *defaultMonitor) IsShutdown() bool {
	return dmt.isShutdown
}

//Shutdown Termination system
func (dmt *defaultMonitor) Shutdown() {
	dmt.isShutdown = true
}

//IncService Add a service
func (dmt *defaultMonitor) IncService() {

}

//DecService Reduce a service
func (dmt *defaultMonitor) DecService() {

}

//WaitService Waiting for all services to end
func (dmt *defaultMonitor) WaitService() {

}

var defMonitor = defaultMonitor{}
