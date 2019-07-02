package monitor

import (
	"magicNet/engine/logger"
	"net/http"
	"strings"
)

// MonitorServer : 监控服务模块 支持https/http协议
type MonitorService struct {
	serviceProtocol string
	serviceTlsCrt   string
	serviceTlsKey   string
	serviceHandle   *http.ServeMux
	servichMethod   http.Handler
}

// Init : 初始化服务
func (M *MonitorService) Init() {
	M.serviceHandle = http.NewServeMux()
	M.serviceProtocol = "http"
	M.serviceTlsCrt = ""
	M.serviceTlsKey = ""
}

// Bind : 绑定服务
func (M *MonitorService) Bind(pattern string, handler http.Handler) {
	M.serviceHandle.Handle(pattern, handler)
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

func (M *MonitorService) lhttp(addr string) bool {
	err := http.ListenAndServe(addr, M.serviceHandle)
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
	err := http.ListenAndServeTLS(addr, M.serviceTlsCrt, M.serviceTlsKey, M.serviceHandle)
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
