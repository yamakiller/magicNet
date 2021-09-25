package listener

import (
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//SpawnWSSListener create an websocket listener
func SpawnWSSListener(l net.Listener, u *websocket.Upgrader) Listener {
	return &WSSListener{_l: l,
		_u: u}
}

//WSSListener Web Socket Listener
type WSSListener struct {
	_l  net.Listener
	_u  *websocket.Upgrader
	_wg sync.WaitGroup
}

//Accept 接受连接
func (slf *WSSListener) Accept(response []interface{}) (interface{}, error) {
	w := response[0].(http.ResponseWriter)
	r := response[1].(*http.Request)

	c, err := slf._u.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	slf._wg.Add(1)

	return &WSSConn{Conn: c, _wg: &slf._wg}, nil
}

//Addr Returns  address
func (slf *WSSListener) Addr() net.Addr {
	return slf._l.Addr()
}

//Close close listener
func (slf *WSSListener) Close() error {
	if err := slf._l.Close(); err != nil {
		return err
	}

	slf._wg.Wait()
	return nil
}

//ToString ....
func (slf *WSSListener) ToString() string {
	return "wss listener"
}

//WSSConn WebSocket connection
type WSSConn struct {
	*websocket.Conn
	_wg *sync.WaitGroup
}

//Close ...
func (slf *WSSConn) Close() error {
	if slf._wg != nil {
		slf._wg.Done()
	}
	slf._wg = nil
	return slf.Conn.Close()
}
