package monitor

/*import (
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"net/http"
	"strings"
)

// MonitorServer : 监控服务模块 支持https/http协议
type MonitorService struct {
	serviceProtocol string
	serviceTlsCrt   string
	serviceTlsKey   string
	serviceMutex   *http.ServeMux
	serviceHandle  *http.Server
	serivceStartMutex util.SpinLock
}

// Init : 初始化服务
func (M *MonitorService) Init() {
	M.serviceMutex = http.NewServeMux()
	M.serviceProtocol = "http"
	M.serviceTlsCrt = ""
	M.serviceTlsKey = ""
}

// Bind : 绑定服务
func (M *MonitorService) Bind(pattern string, handler http.Handler) {
	M.serviceMutex.Handle(pattern, handler)
}

// SetHttps : 设置为Https协议
func (M *MonitorService) SetHttps(tlsCrt, tlsKey string) {
	M.serviceProtocol = "https"
	M.serviceTlsCrt = tlsCrt
	M.serviceTlsKey = tlsKey
}

// Listen ： 启动监听服务
func (M *MonitorService) Listen(addr string) bool {
	if strings.Compare(M.serviceProtocol, "http") == 0 {
		return M.lhttp(addr)
	} else {
		return M.lhttps(addr)
	}
}

func (M *MonitorService) Close() {
	M.serivceStartMutex.Lock()
	defer M.serivceStartMutex.Unlock()

	if M.serviceHandle != nil {
		M.serviceHandle.Close()
		M.serviceHandle = nil
	}
}

func (M *MonitorService) lhttp(addr string) bool {
	M.serivceStartMutex.Lock()
	M.serviceHandle = &http.Server{Addr: addr, Handler: M.serviceMutex}
	M.serivceStartMutex.Unlock()

	err := M.serviceHandle.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info(0, "monitor service closed")
			return true
		} else {
			logger.Error(0, "monitor service start fail:%s", err.Error())
			return false
		}
	}
	return true
}

func (M *MonitorService) lhttps(addr string) bool {
	M.serivceStartMutex.Lock()
	M.serviceHandle = &http.Server{Addr: addr, Handler: M.serviceMutex}
	M.serivceStartMutex.Unlock()

	err := M.serviceHandle.ListenAndServeTLS(M.serviceTlsCrt, M.serviceTlsKey)
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info(0, "monitor service closed")
			return true
		} else {
			logger.Error(0, "monitor service start fail:%s", err.Error())
			return false
		}
	}
	return true
}*/
