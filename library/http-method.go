package library

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/yamakiller/magicNet/engine/util"
)

//HTTPSrvFunc : http 服务函数
type HTTPSrvFunc func(arg1 http.ResponseWriter, arg2 *http.Request)

// IHTTPSrvMethod : HTTP 服务方法接口
type IHTTPSrvMethod interface {
	http.Handler
	RegisterMethod(pattern string, method string, f interface{})
	Close()
	match(requestURI string, method string) interface{}
}

type httpMethodValue struct {
	httpMethod string
	f          interface{}
}

//HTTPSrvMethod : http 服务方法
type HTTPSrvMethod struct {
	suffixRegexp *regexp.Regexp
	methods      map[string]httpMethodValue
	l            sync.RWMutex
}

//NewHTTPSrvMethod 新建一个服务方法
func NewHTTPSrvMethod() IHTTPSrvMethod {
	r := &HTTPSrvMethod{}
	r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
	r.methods = make(map[string]httpMethodValue, 32)
	return r
}

func (hsm *HTTPSrvMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hh-", r.Method)
	f := hsm.match(r.RequestURI, r.Method)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if v, ok := f.(HTTPSrvFunc); ok {
		v(w, r)
	}
}

// RegisterMethod : 注册服务方法
func (hsm *HTTPSrvMethod) RegisterMethod(pattern string, method string, f interface{}) {
	hsm.l.Lock()
	defer hsm.l.Unlock()
	hsm.methods[pattern] = httpMethodValue{httpMethod: method, f: f}
}

// Close : 关闭方法服务
func (hsm *HTTPSrvMethod) Close() {
	hsm.l.Lock()
	defer hsm.l.Unlock()
	hsm.methods = make(map[string]httpMethodValue)
}

func (hsm *HTTPSrvMethod) match(requestURI string, method string) interface{} {
	hsm.l.RLock()
	defer hsm.l.RUnlock()
	suffix := hsm.suffixRegexp.FindStringSubmatch(requestURI)
	if suffix != nil && len(suffix) >= 1 {
		return nil
	}

	tmpURI := requestURI
	idx := strings.LastIndex(requestURI, "?")
	if idx > 0 {
		tmpURI = util.SubStr2(tmpURI, 0, idx)
	}

	r := hsm.methods[tmpURI]

	if strings.Index(strings.ToLower(r.httpMethod),
		strings.ToLower(method)) >= 0 {
		return r.f
	}

	return nil
}
