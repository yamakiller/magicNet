package service

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"net/http"
	"strings"
)

// MonitorService : 监视去服务
type MonitorService struct {
	Service
	Proto      string
	Addr       string
	httpMutex  *http.ServeMux
	httpHandle *http.Server
}

// Started : 监视器启动函数
func (ms *MonitorService) Started(context actor.Context) {
	//1.读取配置文件信息
	//2.启动HTTP/HTTPS服务

	ms.Service.Started(context)
}

// Stoped : 停止服务
func (ms *MonitorService) Stoped(context actor.Context) {
	ms.httpHandle.Close()
}

// 启动监听 addr 格式 ip:port
func (ms *MonitorService) listen() bool {
	ms.httpHandle = &http.Server{Addr: ms.Addr, Handler: ms.httpMutex}
	if strings.Compare(ms.Proto, "http") == 0 {
		return ms.lhttp()
	}
	return ms.lhttps()
}

func (ms *MonitorService) lhttp() bool {

	err := ms.httpHandle.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info(0, "monitor service closed")
			return true
		}
		logger.Error(0, "monitor service start fail:%s", err.Error())
		return false
	}
	return true
}

func (ms *MonitorService) lhttps() bool {

	err := ms.httpHandle.ListenAndServeTLS("", "")
	if err != nil {
		if err == http.ErrServerClosed {
			logger.Info(0, "monitor service closed")
			return true
		}
		logger.Error(0, "monitor service start fail:%s", err.Error())
		return false
	}
	return true
}
