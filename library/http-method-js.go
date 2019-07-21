package library

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/yamakiller/magicNet/engine/files"
	"github.com/yamakiller/magicNet/script"
	"github.com/yamakiller/magicNet/script/stack"

	"github.com/robertkrimen/otto"
)

// HTTPSrvMethodJS : js服务解析方法
type HTTPSrvMethodJS struct {
	HTTPSrvMethod
	curHTTPRequest *http.Request
	jsChan         chan int
	jsStop         chan int
	jsVim          *stack.JSStack
}

//NewHTTPSrvMethodJS 新建一个带JS功能的服务方法
func NewHTTPSrvMethodJS() IHTTPSrvMethod {
	r := &HTTPSrvMethodJS{}
	r.suffixRegexp, _ = regexp.Compile(`\.\w+.*`)
	r.methods = make(map[string]httpMethodValue, 32)
	r.jsChan = make(chan int, 1)
	r.jsStop = make(chan int)
	r.jsVim = script.NewJSStack()

	r.jsVim.SetFunc("GetHttpMethod", r.getHTTPMethodJS)
	r.jsVim.SetFunc("GetHttpParam", r.getHTTPParamJS)
	return r
}

func (hsm *HTTPSrvMethodJS) getHTTPMethodJS(js otto.FunctionCall) otto.Value {
	vmst, err := otto.ToValue(hsm.curHTTPRequest.Method)
	if err != nil {
		panic(err)
	}

	return vmst
}

func (hsm *HTTPSrvMethodJS) getHTTPParamJS(js otto.FunctionCall) otto.Value {
	r := hsm.curHTTPRequest
	m := strings.ToLower(r.Method)
	if m == "post" {
		result, _ := js.Otto.Object(`({})`)
		for k, v := range r.PostForm {
			result.Set(k, v)
		}

		vmst, err := otto.ToValue(result)
		if err != nil {
			panic(err)
		}

		return vmst
	}

	i := 0
	result, _ := js.Otto.Object(`({})`)
	for k, v := range r.Form {
		result.Set(k, v)
		i++
	}

	if i == 0 {
		for k, v := range r.URL.Query() {
			result.Set(k, v)
			i++
		}
	}

	vmst, err := otto.ToValue(result)
	if err != nil {
		panic(err)
	}

	return vmst
}

func (hsm *HTTPSrvMethodJS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f := hsm.match(r.RequestURI, r.Method)
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if v, ok := f.(string); ok {
		hsm.runJs(v, w, r)
		return
	}

	if v, ok := f.(func(http.ResponseWriter, *http.Request)); ok {
		v(w, r)
	}
}

// Close : 关闭方法注册服务
func (hsm *HTTPSrvMethodJS) Close() {
	hsm.l.Lock()
	defer hsm.l.Unlock()
	hsm.methods = make(map[string]httpMethodValue)
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
		hsm.curHTTPRequest = nil
		select {
		case <-hsm.jsStop:
			return
		case <-hsm.jsChan:
		}
	}()

	hsm.curHTTPRequest = r
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
