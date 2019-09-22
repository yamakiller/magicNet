package implement

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/library"
	"github.com/yamakiller/magicNet/service"
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
	service.Service
	Proto       string
	Addr        string
	OAuto2      *library.OAuth2
	MakerMethod MakeHTTPMethod
	CertFile    string
	KeyFile     string
	Regiter     func(pid *actor.PID, m library.IHTTPSrvMethod)

	isShutdown bool
	httpErr    error
	httpMethod library.IHTTPSrvMethod
	httpWait   sync.WaitGroup
	httpMutex  *http.ServeMux
	httpHandle *http.Server
}

// Init : 初始化服务
func (ms *MonitorService) Init() {
	ms.Service.Init()
	ms.RegisterMethod(&actor.Started{}, ms.Started)
	ms.RegisterMethod(&actor.Stopped{}, ms.Stoped)
}

// Started : 监视器启动函数
func (ms *MonitorService) Started(context actor.Context, message interface{}) {
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
		logger.Info(context.Self().ID, "OAuto2 Config Auth-token-exp:%d Sec", ms.OAuto2.TokenExp)
		logger.Info(context.Self().ID, "OAuto2 Config Auth-refresh-token-exp:%d Sec", ms.OAuto2.RefreshTokenExp)
		logger.Info(context.Self().ID, "OAuto2 Config Auth-is-generate-refresh-token:%v Sec", ms.OAuto2.IsGenerateRefresh)
		logger.Info(context.Self().ID, "OAuto2 Config S256key:%s", ms.OAuto2.S256Key)
		logger.Info(context.Self().ID, "OAuto2 Config Access-token-url:%s", ms.OAuto2.AccessURI)
	}

	if ms.Regiter != nil {
		ms.Regiter(context.Self(), ms.httpMethod)
	}

	ln, err := ms.listen()
	ms.httpErr = err
	if ms.httpErr != nil {
		logger.Error(context.Self().ID, "%s %s service start fail:%v", ms.Name(), ms.Proto, ms.httpErr)
		goto end_lable
	}

	logger.Info(context.Self().ID, "%s %s service start success[addr:%s]", ms.Name(), ms.Proto, ms.Addr)
	if err == nil {
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
						ms.httpErr = ms.httpHandle.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, ms.CertFile, ms.KeyFile)
					}
				} else {
					time.Sleep(time.Millisecond * 10)
				}
			}

			ms.httpWait.Done()
		}()
	}
end_lable:
	ms.Service.Started(context, message)
}

// Stoped : 停止服务
func (ms *MonitorService) Stoped(context actor.Context, message interface{}) {
	err := ms.httpHandle.Close()
	if err != http.ErrServerClosed {
		logger.Warning(context.Self().ID, "monitor service close error:%v", err)
	}
	ms.Service.Stoped(context, message)
	//!位置可以考虑一下
	ms.httpMethod.Close()
}

// Shutdown 关闭服务
func (ms *MonitorService) Shutdown() {
	ms.isShutdown = true
	ms.Service.Shutdown()
	ms.httpWait.Wait()
}

// 启动监听 addr 格式 ip:port
func (ms *MonitorService) listen() (net.Listener, error) {
	addr := ms.httpHandle.Addr
	if addr == "" {
		addr = ":" + ms.Proto
	}

	return net.Listen("tcp", addr)
}

var (
	_ service.IService = &MonitorService{}
)
