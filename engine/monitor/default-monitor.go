package monitor

//DefaultMonitor : 默认监控器
type defaultMonitor struct {
	isShutdown bool
}

// IsShutdown : 系统是否已经关闭
func (dmt *defaultMonitor) IsShutdown() bool {
	return dmt.isShutdown
}

// Shutdown : 关闭系统
func (dmt *defaultMonitor) Shutdown() {
	dmt.isShutdown = true
}

// IncService : 增加一个服务
func (dmt *defaultMonitor) IncService() {

}

// DecService : 减少一个服务
func (dmt *defaultMonitor) DecService() {

}

// WaitService : 等待所有服务结束
func (dmt *defaultMonitor) WaitService() {

}

var defMonitor = defaultMonitor{}
