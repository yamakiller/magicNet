package monitor

import (
	"sync"
	"sync/atomic"
)

//DefaultMonitor : 监视器
type DefaultMonitor struct {
	waitValue  sync.WaitGroup
	isShutdown int32
}

// IsShutdown : 是否已经关闭
func (dm *DefaultMonitor) IsShutdown() bool {
	if atomic.LoadInt32(&dm.isShutdown) != 0 {
		return true
	}
	return false
}

// Shutdown : 关闭程序
func (dm *DefaultMonitor) Shutdown() {
	atomic.CompareAndSwapInt32(&dm.isShutdown, 0, 1)
}

// IncService : 增加一个服务
func (dm *DefaultMonitor) IncService() {
	dm.waitValue.Add(1)
}

// DecService : 减少一个服务
func (dm *DefaultMonitor) DecService() {
	dm.waitValue.Done()
}

// WaitService : 等待所有服务结束
func (dm *DefaultMonitor) WaitService() {
	dm.waitValue.Wait()
}
