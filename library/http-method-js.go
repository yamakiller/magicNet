package library

import (
	"magicNet/engine/files"
	"net/http"
	"regexp"
	"strings"
)

// HTTPSrvMethodJS : js服务解析方法
type HTTPSrvMethodJS struct {
	HTTPSrvMethod
}

//NewHTTPSrvMethodJS 新建一个带JS功能的服务方法
func NewHTTPSrvMethodJS() IHTTPSrvMethod {
	r := &HTTPSrvMethodJS{}
	r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
	r.methods = make(map[string]interface{}, 32)
	return r
}

func (hsm *HTTPSrvMethodJS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := hsm.match(r.RequestURI)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if v, ok := f.(string); ok {
		filNme, filFun, err := hsm.decomposeJs(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		hsm.runJs(filNme, filFun, w, r)
	} else if v, ok := f.(HTTPSrvFunc); ok {
		v(w, r)
	}
}

//RegisterMethod : 注册
func (hsm *HTTPSrvMethodJS) RegisterMethod(pattern string, f interface{}) {
	hsm.l.Lock()
	defer hsm.l.Unlock()
	hsm.methods[pattern] = f
}

func (hsm *HTTPSrvMethodJS) runJs(jsfile string, jsfun string, w http.ResponseWriter, r *http.Request) {
	fileFullPath := files.GetFullPathForFilename(jsfile)
	if !files.IsFileExist(fileFullPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := files.GetDataFromFile(fileFullPath)
	if data.IsNil() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	//虚拟机器，获取结果
}

func (hsm *HTTPSrvMethodJS) decomposeJs(v string) (string, string, error) {
	s := strings.Split(v, "#")
	if len(s) != 2 {
		return "", "", http.ErrNotSupported
	}

	return s[0], s[1], nil
}
