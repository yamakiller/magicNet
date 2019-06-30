package monitor

import (
	"magicNet/logger"
	"magicNet/util"
	"strings"
	"sync"
	"sync/atomic"
)

// Monitor 监视器
type Monitor struct {
	w sync.WaitGroup
	a int32          /*Actor number*/
	e bool           /*is shutdown ?*/
	s string         /*system run state */
	h MonitorService /*system monitor http service*/
}

const (
	monitorIdle     = "idle"
	monitorStart    = "starting"
	monitorRun      = "running"
	monitorShutdown = "shutdown"
)

var instMonitor *Monitor

// Init : 初始化监视器
func Init() {
	instMonitor = &Monitor{sync.WaitGroup{}, 0, false, monitorIdle, MonitorService{}}

}

// StartService : 启动服务
func StartService() bool {
	msc := util.GetEnvInstance().GetMap("monitor")
	if msc == nil {
		return true
	}

	logger.Info(0, "monitor service starting...")
	instMonitor.h.Init()
	instMonitor.h.Bind("/", &MonitorMethod{})

	protocol := msc["protocol"].String()
	address := msc["address"].String()
	port := msc["port"].String()
	if strings.Compare(protocol, "https") == 0 {
		instMonitor.h.SetHttps(msc["tls-crt"].String(),
			msc["tls-key"].String())
	}

	WaitInc()
	go func(addr string, proto string) {
		logger.Info(0, "monitor service %s %s", proto, addr)
		if !instMonitor.h.Listen(addr) {
			instMonitor.e = true
		}
		WaitDec()
	}(address+":"+port, protocol)

	return true
}

// WaitDec : 完成一个等待
func WaitDec() {
	instMonitor.w.Done()
}

// WaitInc : 等待+1
func WaitInc() {
	instMonitor.w.Add(1)
}

// WaitSupper : 等待系统正确退出
func WaitSupper() {
	instMonitor.w.Wait()
}

// ActorInc : Actor 数量增 1
func ActorInc() {
	atomic.AddInt32(&instMonitor.a, 1)
}

// ActorDec : Actor 数量减 1
func ActorDec() {
	atomic.AddInt32(&instMonitor.a, -1)
}

// ActorCount : 获取当前 Acotr 的总数
func ActorCount() int {
	return int(atomic.LoadInt32(&instMonitor.a))
}

// SetStateIdle : 设置系统为闲置状态
func SetStateIdle() {
	instMonitor.s = monitorIdle
}

// SetStateStart : 设置系统为启动中
func SetStateStart() {
	instMonitor.s = monitorStart
}

// SetStateRun : 设置系统为运行中
func SetStateRun() {
	instMonitor.s = monitorRun
}

// SetStateShutdown : 设置系统为终止中
func SetStateShutdown() {
	instMonitor.s = monitorShutdown
}

// Shutdown : 终止系统运行
func Shutdown() {
	instMonitor.e = true
}

// IsShutdown : 系统是否已终止
func IsShutdown() bool {
	return instMonitor.e
}
