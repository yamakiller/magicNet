package implement

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/library"
	"github.com/yamakiller/magicNet/logger"
	"github.com/yamakiller/magicNet/service"
)

//tcpKeepAliveListener : Overload net/http tcpKeepAliveListener
type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept : Overload net/http tcpKeepAliveListener.Accept
func (slf tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := slf.AcceptTCP()
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

// Init : Initialization service
func (slf *MonitorService) Init() {
	slf.Service.Init()
	slf.RegisterMethod(&actor.Started{}, slf.Started)
	slf.RegisterMethod(&actor.Stopping{}, slf.Stopping)
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
}

// Started : Monitor startup function
func (slf *MonitorService) Started(context actor.Context,
	sender *actor.PID,
	message interface{}) {

	slf.isShutdown = false
	slf.httpMutex = http.NewServeMux()
	slf.httpHandle = &http.Server{Addr: slf.Addr, Handler: slf.httpMutex}
	if slf.MakerMethod == nil {
		slf.httpMethod = library.NewHTTPSrvMethod()
	} else {
		slf.httpMethod = slf.MakerMethod()
	}
	slf.httpMutex.Handle("/", slf.httpMethod)

	if slf.OAuto2 != nil {
		slf.OAuto2.Init(slf.httpMethod)
		logger.Info(context.Self().ID, "OAuto2 Config Auth-token-exp:%d Sec", slf.OAuto2.TokenExp)
		logger.Info(context.Self().ID, "OAuto2 Config Auth-refresh-token-exp:%d Sec", slf.OAuto2.RefreshTokenExp)
		logger.Info(context.Self().ID, "OAuto2 Config Auth-is-generate-refresh-token:%v Sec", slf.OAuto2.IsGenerateRefresh)
		logger.Info(context.Self().ID, "OAuto2 Config S256key:%s", slf.OAuto2.S256Key)
		logger.Info(context.Self().ID, "OAuto2 Config Access-token-url:%s", slf.OAuto2.AccessURI)
	}

	if slf.Regiter != nil {
		slf.Regiter(context.Self(), slf.httpMethod)
	}

	ln, err := slf.listen()
	slf.httpErr = err
	if slf.httpErr != nil {
		logger.Error(context.Self().ID, "%s %s service start fail:%v", slf.Name(), slf.Proto, slf.httpErr)
		goto end_lable
	}

	logger.Info(context.Self().ID, "%s %s service start success[addr:%s]", slf.Name(), slf.Proto, slf.Addr)
	if err == nil {
		slf.httpWait.Add(1)
		go func() {
			for {
				if slf.isShutdown {
					break
				}

				if slf.httpErr == nil {
					if slf.Proto == "http" {
						slf.httpErr = slf.httpHandle.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
					} else {
						defer ln.Close()
						slf.httpErr = slf.httpHandle.ServeTLS(tcpKeepAliveListener{ln.(*net.TCPListener)}, slf.CertFile, slf.KeyFile)
					}
				} else {
					time.Sleep(time.Millisecond * 10)
				}
			}

			slf.httpWait.Done()
		}()
	}
end_lable:
	slf.Service.Started(context, sender, message)
}

// Stopping : Out of service
func (slf *MonitorService) Stopping(context actor.Context, sender *actor.PID, message interface{}) {
	err := slf.httpHandle.Close()
	if err != http.ErrServerClosed {
		logger.Warning(context.Self().ID, "monitor service close error:%v", err)
	}

	slf.httpMethod.Close()
}

//Stoped Service has stopped
func (slf *MonitorService) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.Service.Stoped(context, sender, message)
}

// Shutdown Close service
func (slf *MonitorService) Shutdown() {
	slf.isShutdown = true
	slf.Service.Shutdown()
	slf.httpWait.Wait()
}

// Start monitoring addr format ip:port
func (slf *MonitorService) listen() (net.Listener, error) {
	addr := slf.httpHandle.Addr
	if addr == "" {
		addr = ":" + slf.Proto
	}

	return net.Listen("tcp", addr)
}

var (
	_ service.IService = &MonitorService{}
)
