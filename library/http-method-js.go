package library

import (
	"fmt"
	"magicNet/engine/files"
	"magicNet/script"
	"magicNet/script/stack"
	"net/http"
	"regexp"
)

// HTTPSrvMethodJS : js服务解析方法
type HTTPSrvMethodJS struct {
	HTTPSrvMethod
	jsChan chan int
	jsStop chan int
	jsVim  *stack.JSStack
}

//NewHTTPSrvMethodJS 新建一个带JS功能的服务方法
func NewHTTPSrvMethodJS() IHTTPSrvMethod {
	r := &HTTPSrvMethodJS{}
	r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
	r.methods = make(map[string]interface{}, 32)
	r.jsChan = make(chan int, 1)
	r.jsStop = make(chan int)
	r.jsVim = script.NewJSStack()
	return r
}

func (hsm *HTTPSrvMethodJS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := hsm.match(r.RequestURI)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if v, ok := f.(string); ok {

		hsm.runJs(v, w, r)
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

// Close : 关闭方法注册服务
func (hsm *HTTPSrvMethodJS) Close() {
	hsm.l.Lock()
	defer hsm.l.Unlock()
	hsm.methods = make(map[string]interface{})
	close(hsm.jsStop)
	close(hsm.jsChan)
	//? 需要吗
	hsm.jsVim = nil
}

func (hsm *HTTPSrvMethodJS) runJs(jsfile string, w http.ResponseWriter, r *http.Request) {
	fileFullPath := files.GetFullPathForFilename(jsfile)
	if !files.IsFileExist(fileFullPath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	select {
	case <-hsm.jsStop:
		return
	case hsm.jsChan <- 1:
	}

	defer func() {
		select {
		case <-hsm.jsStop:
			return
		case <-hsm.jsChan:
		}
	}()

	result, err := hsm.jsVim.ExecuteScriptFile(jsfile)
	if err != nil {
		if err == stack.ErrJSNotFindFile ||
			err == stack.ErrJSNotFileData {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg := fmt.Sprintf("{code: 140, message:'%s'}", err.Error())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
		return
	}

	msg, err := result.ToString()
	w.WriteHeader(http.StatusOK)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{code: 141, message:'%s'}", err.Error())))
		return
	}
	w.Write([]byte(msg))
}
