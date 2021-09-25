package borker

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yamakiller/magicLibs/net/listener"
)

//WSSBorker websocket borker server
type WSSBorker struct {
	WSPath           string
	Spawn            func(*listener.WSSConn) error
	ReadBufferSize   int
	WriteBufferSize  int
	HandshakeTimeout time.Duration
	_listen          listener.Listener
	_httpHandler     *http.ServeMux
	_httpListen      net.Listener
	_wg              sync.WaitGroup
}

//ListenAndServe 监听并启动服务
func (slf *WSSBorker) ListenAndServe(addr string) error {
	if slf.WSPath == "" {
		slf.WSPath = "/"
	}
	if strings.Index(slf.WSPath, "/") < 0 {
		slf.WSPath = fmt.Sprintf("/%s", slf.WSPath)
	}

	slf._httpHandler = http.NewServeMux()
	slf._httpHandler.HandleFunc(slf.WSPath, slf.onWSS)

	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	lst, err := net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}
	slf._httpListen = lst
	slf._listen = listener.SpawnWSSListener(lst, &websocket.Upgrader{
		ReadBufferSize:   slf.ReadBufferSize,
		WriteBufferSize:  slf.WriteBufferSize,
		HandshakeTimeout: slf.HandshakeTimeout,
	})

	slf._wg.Add(1)
	go slf.Serve()

	return nil
}

func (slf *WSSBorker) ListenAndServeTls(addr string, ptls *tls.Config) error {
	return errors.New("undefined tls")
}

//Serve accept connection
func (slf *WSSBorker) Serve() {
	defer slf._wg.Done()
	http.Serve(slf._httpListen, slf._httpHandler)
}

func (slf *WSSBorker) onWSS(w http.ResponseWriter, r *http.Request) {
	params := make([]interface{}, 2)
	params[0] = w
	params[1] = r
	c, err := slf._listen.Accept(params)
	if err != nil {
		return
	}

	if e := slf.Spawn(c.(*listener.WSSConn)); e != nil {
		c.(*listener.WSSConn).Close()
	}
}

//Listener Returns listenner object
func (slf *WSSBorker) Listener() listener.Listener {
	return slf._listen
}

//Shutdown 关闭服务
func (slf *WSSBorker) Shutdown() {
	slf._httpListen.Close()
	slf._listen.Close()
	slf._wg.Wait()
	slf._listen = nil
}
