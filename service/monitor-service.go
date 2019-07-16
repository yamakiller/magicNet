package service

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/library"
	"net"
	"net/http"
	"sync"
	"time"
)

//tcpKeepAliveListener : 重载net/http tcpKeepAliveListener
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept : 重载net/http tcpKeepAliveListener.Accept
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// MakeHTTPMethod : http method 生成器
type MakeHTTPMethod func() library.IHTTPSrvMethod

// MonitorService : 监视去服务
type MonitorService struct {
	Service
	Proto       string
	Addr        string
	OAuto2      *library.OAuth2
	MakerMethod MakeHTTPMethod
	certFile    string
	keyFile     string

	isShutdown bool
	httpErr    error
	httpMethod library.IHTTPSrvMethod
	httpWait   sync.WaitGroup
	httpMutex  *http.ServeMux
	httpHandle *http.Server
}

// Started : 监视器启动函数
func (ms *MonitorService) Started(context actor.Context) {
	ms.isShutdown = false
	ms.httpMutex = http.NewServeMux()
	ms.httpHandle = &http.Server{Addr: ms.Addr, Handler: ms.httpMutex}
	if ms.MakerMethod == nil {
		ms.httpMethod = library.NewHTTPSrvMethod()
	} else {
		ms.httpMethod = ms.MakerMethod()
	}
	ms.httpMutex.Handle("/", ms.httpMethod)

	if ms.OAuto2 != nil {
		ms.OAuto2.Init(ms.httpMethod)
	}

	ln, err := ms.listen()
	ms.httpErr = err
	if ms.httpErr != nil {
		logger.Error(context.Self().ID, "monitor %s service start fail:%v", ms.Proto, ms.httpErr)
		goto end_lable
	}

	logger.Info(context.Self().ID, "monitor %s service start success[addr:%s]", ms.Proto, ms.Addr)
	if err != nil {
		ms.httpWait.Add(1)
		go func() {

			for {
				if ms.isShutdown {
					break
				}

				if ms.httpErr == nil {
					if ms.Proto == "http" {
						ms.httpErr = ms.httpHandle.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
					} else {
						defer ln.Close()
						ms.httpErr = ms.httpHandle.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, ms.certFile, ms.keyFile)
					}
				} else {
					time.Sleep(time.Millisecond * 10)
				}
			}

			ms.httpWait.Done()
		}()
	}
end_lable:
	ms.Service.Started(context)
}

// Stoped : 停止服务
func (ms *MonitorService) Stoped(context actor.Context) {
	err := ms.httpHandle.Close()
	if err != http.ErrServerClosed {
		logger.Warning(context.Self().ID, "monitor service close error:%v", err)
	}
	ms.Service.Stoped(context)
}

// Shutdown 关闭服务
func (ms *MonitorService) Shutdown() {
	if ms.pid == nil {
		return
	}

	ms.isShutdown = true
	ms.pid.Stop()
	ms.httpWait.Wait()
	ms.wait.Wait()
}

// 启动监听 addr 格式 ip:port
func (ms *MonitorService) listen() (net.Listener, error) {
	addr := ms.httpHandle.Addr
	if addr == "" {
		addr = ":" + ms.Proto
	}

	return net.Listen("tcp", addr)
}
