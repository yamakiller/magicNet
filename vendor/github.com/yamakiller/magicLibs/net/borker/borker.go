package borker

import (
	"crypto/tls"

	"github.com/yamakiller/magicLibs/net/listener"
)

//Borker 网络代理服务
type Borker interface {
	ListenAndServe(string) error
	ListenAndServeTls(string, *tls.Config) error
	Listener() listener.Listener
	Serve() error
	Shutdown()
}
