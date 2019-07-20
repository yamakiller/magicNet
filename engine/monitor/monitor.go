package monitor

// Monitor : 监视器模型接口
type Monitor interface {
	// IsShutdown : 系统是否已经关闭
	IsShutdown() bool

	// Shutdown : 关闭系统
	Shutdown()

	// IncService : 增加一个服务
	IncService()

	// DecService : 减少一个服务
	DecService()

	// WaitService : 等待所有服务结束
	WaitService()
}

var systemMonitor Monitor

// WithMonitor 关联到系统监视器
func WithMonitor(m Monitor) {
	systemMonitor = m
}

// IsShutdown : 系统是否已被终止
func IsShutdown() bool {
	if systemMonitor == nil {
		return defMonitor.IsShutdown()
	}
	return systemMonitor.IsShutdown()
}

// Shutdown : 终止系统
func Shutdown() {
	if systemMonitor == nil {
		defMonitor.Shutdown()
		return
	}
	systemMonitor.Shutdown()
}

// IncService : 增加一个服务
func IncService() {
	if systemMonitor == nil {
		defMonitor.IncService()
		return
	}
	systemMonitor.IncService()
}

// DecService : 减少一个服务
func DecService() {
	if systemMonitor == nil {
		defMonitor.DecService()
		return
	}
	systemMonitor.DecService()
}

// WaitService : 等待所有服务结束
func WaitService() {
	if systemMonitor == nil {
		defMonitor.WaitService()
		return
	}
	systemMonitor.WaitService()
}
