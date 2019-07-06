package monitor

import (
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"magicNet/engine/hook"
	"strings"
	//"reflect"
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
	hmethod *MonitorMethod /*system monitor http service method*/
}

type monitorConfig struct {
	protocol string
	address  string
	port     string
}

const (
	monitorIdle     = "idle"
	monitorStart    = "starting"
	monitorRun      = "running"
	monitorShutdown = "shutdown"
)

var instMonitor *Monitor
var monitorInitHook hook.InitializeHook

// Init : 初始化监视器
func Init() {
	instMonitor = &Monitor{sync.WaitGroup{}, 0, false, monitorIdle, MonitorService{}, NewMonitorMethod()}
}

// SetMonitorInitHook : 设置监视器初始/销毁Hook函数
func SetMonitorInitHook(miHk hook.InitializeHook) {
	if monitorInitHook == nil {
		monitorInitHook = miHk
	}
}

// StartService : 启动服务
func StartService() bool {

	msc := util.GetEnvMap(util.GetEnvRoot(), "monitor")
	if msc == nil {
		return true
	}

	logger.Info(0, "monitor service starting")
	instMonitor.h.Init()
	instMonitor.h.Bind("/", instMonitor.hmethod)

	protocol := util.GetEnvString(msc, "protocol", "http")
	address := util.GetEnvString(msc, "address", "127.0.0.1")
	port := util.GetEnvString(msc, "port", "8001")
	if strings.Compare(protocol, "https") == 0 {
		instMonitor.h.SetHttps(util.GetEnvString(msc, "tls-crt", ""),
													 util.GetEnvString(msc, "tls-key", ""))
	}

	if (!monitorInitHook.Initialize()) {
		return false
	}

	WaitInc()
	go func(addr string, proto string) {
		logger.Info(0, "monitor service %s %s", proto, addr)
		if !instMonitor.h.Listen(addr) {
			instMonitor.e = true
		}
		monitorInitHook.Finalize()
		WaitDec()
	}(address+":"+port, protocol)

	return true
}

// StopService : 停止服务
func StopService() {
	instMonitor.h.Close()
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

// RegisterHttpMethod : 註冊HTTP 映射方法
func RegisterHttpMethod(pattern string, f MonitorHttpFunction) {
	instMonitor.hmethod.methods[pattern] = f
}

// Shutdown : 终止系统运行
func Shutdown() {
	instMonitor.e = true
}

// IsShutdown : 系统是否已终止
func IsShutdown() bool {
	return instMonitor.e
}