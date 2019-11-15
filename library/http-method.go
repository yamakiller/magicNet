package library

import (
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/yamakiller/magicLibs/util"
)

//HTTPSrvFunc : Http service function
type HTTPSrvFunc func(arg1 http.ResponseWriter, arg2 *http.Request)

// IHTTPSrvMethod : HTTP service method interface
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

//HTTPSrvMethod : Http service method
type HTTPSrvMethod struct {
	suffixRegexp *regexp.Regexp
	methods      map[string]httpMethodValue
	l            sync.RWMutex
}

//NewHTTPSrvMethod Create a new service method
func NewHTTPSrvMethod() IHTTPSrvMethod {
	r := &HTTPSrvMethod{}
	r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
	r.methods = make(map[string]httpMethodValue, 32)
	return r
}

func (slf *HTTPSrvMethod) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := slf.match(r.RequestURI, r.Method)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if v, ok := f.(HTTPSrvFunc); ok {
		v(w, r)
	}
}

// RegisterMethod : Registration service method
func (slf *HTTPSrvMethod) RegisterMethod(pattern string, method string, f interface{}) {
	slf.l.Lock()
	defer slf.l.Unlock()
	slf.methods[pattern] = httpMethodValue{httpMethod: method, f: f}
}

// Close : Close method service
func (slf *HTTPSrvMethod) Close() {
	slf.l.Lock()
	defer slf.l.Unlock()
	slf.methods = make(map[string]httpMethodValue)
}

func (slf *HTTPSrvMethod) match(requestURI string, method string) interface{} {
	slf.l.RLock()
	defer slf.l.RUnlock()
	suffix := slf.suffixRegexp.FindStringSubmatch(requestURI)
	if suffix != nil && len(suffix) >= 1 {
		return nil
	}

	tmpURI := requestURI
	idx := strings.LastIndex(requestURI, "?")
	if idx > 0 {
		tmpURI = util.SubStr2(tmpURI, 0, idx)
	}

	r := slf.methods[tmpURI]

	if strings.Index(strings.ToLower(r.httpMethod),
		strings.ToLower(method)) >= 0 {
		return r.f
	}

	return nil
}
