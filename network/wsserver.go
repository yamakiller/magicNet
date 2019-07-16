package network

import (
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"magicNet/engine/util"
	"magicNet/timer"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// wsTCPKeepAliveListener : 重载net/http tcpKeepAliveListener
type wsTCPKeepAliveListener struct {
	*net.TCPListener
}

// Accept : 重载net/http wsTCPKeepAliveListener.Accept
func (ln wsTCPKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// WSListen : WebSocket 监听器
type wslisten struct {
	handle     int32
	accept     websocket.Upgrader
	conns      sync.Map
	httpSrv    *http.Server
	httpMtx    *http.ServeMux
	httpErr    error
	httpWait   sync.WaitGroup
	isShutdown bool
}

// Listen :
func (wsl *wslisten) listen(operator *actor.PID, addr string) error {
	wsl.accept = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsl.httpMtx = http.NewServeMux()
	wsl.httpSrv = &http.Server{Addr: addr, Handler: wsl.httpMtx}
	wsl.httpMtx.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s, err := wsl.accept.Upgrade(w, r, nil)
		if err != nil {
			//错误日志
			logger.Fatal(operator.ID, "web socket accept fail:%v", err)
			return
		}

		handle, so := operGrap()
		if handle == -1 || so == nil {
			s.Close()
			logger.Error(operator.ID, "web socket accept error:%s", "lack of socket resources")
			return
		}

		now := timer.Now()
		so.l.Lock()
		if wsl.isShutdown {
			s.Close()
			so.b = resIdle
			so.l.Unlock()
			return
		}

		so.s = &wsconn{h: handle, s: s, o: operator, stat: Connecting, out: make(chan *NetChunk, wsOutChanMax)}
		conn, _ := so.s.(*wsconn)
		conn.i.ReadLastTime = now
		conn.i.WriteLastTime = now

		conn.w.Add(1)
		go func(c *wsconn, cso *slot) {
			for {
				if c.stat != Connecting && c.stat != Connected {
					goto read_end
				}

				msgType, data, err := c.s.ReadMessage()
				if err != nil {
					//记录错误日志
					goto read_error
				}

				// 不接收非二进制编码数据
				if msgType != websocket.BinaryMessage {
					goto read_error
				}

				// 丢弃数据包
				if c.stat != Connected {
					continue
				}

				c.i.ReadBytes += uint64(len(data))
				c.i.ReadLastTime = timer.Now()
				//数据包丢给 Actor
				actor.DefaultSchedulerContext.Send(c.o, &NetChunk{Data: data})
			}
		read_error:
			c.stat = Closing
			c.s.Close()
		read_end:
			var (
				closeHandle   int32
				closeOperator *actor.PID
			)

			cso.l.Lock()
			closeHandle = c.h
			closeOperator = c.o
			close(c.out)
			//-----等待写协程结束------
			for {
				if atomic.CompareAndSwapInt32(&c.outStat, 1, 1) {
					break
				}
			}
			//----------------------
			wsl.conns.Delete(c.h)

			cso.s = nil
			cso.b = resIdle
			cso.l.Unlock()

			c.w.Done()

			actor.DefaultSchedulerContext.Send(closeOperator, &NetClose{handle: closeHandle})

		}(conn, so)

		conn.w.Add(1)
		go func(c *wsconn) {
			for {
				if c.stat != Connecting && c.stat != Connected {
					goto write_end
				}

				select {
				case msg := <-c.out:
					if c.stat != Connecting && c.stat != Connected {
						goto write_end
					}

					if err := c.s.WriteMessage(websocket.BinaryMessage, msg.Data); err != nil {
						goto write_error
					}

					c.i.WriteBytes += uint64(len(msg.Data))
					c.i.WriteLastTime = timer.Now()
				}
			}
		write_error:
			c.stat = Closing
		write_end:
			c.w.Done()
			c.outStat = 1
		}(conn)

		so.b = resAssigned
		so.l.Unlock()

		wsl.conns.Store(handle, int32(1))
	})

	ln, err := wsl.httpListen(addr)
	if err != nil {
		return err
	}

	wsl.httpWait.Add(1)
	go func() {
		for {
			wsl.httpErr = wsl.httpSrv.Serve(wsTCPKeepAliveListener{ln.(*net.TCPListener)})
			wsl.isShutdown = true
			break
		}

		wsl.httpWait.Done()
	}()
	// 启动闲置检测器
	wsl.httpWait.Add(1)
	go func() {
		for {
			if wsl.isShutdown {
				break
			}

			time.Sleep(time.Second * 1)

			now := timer.Now()
			wsl.conns.Range(func(handle interface{}, v interface{}) bool {
				so := operGet(handle.(int32))
				if so == nil {
					return true
				}

				if so.b == resIdle {
					return true
				}

				so.l.Lock()
				defer so.l.Unlock()
				if so.b == resIdle || so.b == resOccupy {
					return true
				}

				// 维护KeepAlive
				if conn, ok := so.s.(*wsconn); ok {
					if conn.keepAive == 0 {
						return true
					}

					if (now - conn.i.ReadLastTime) > conn.keepAive {
						conn.close(nil)
					}
				}
				return true
			})
		}

		//------------------关闭所有连接-----------------------------
		wsl.conns.Range(func(handle interface{}, v interface{}) bool {
			so := operGet(handle.(int32))
			if so.b == resIdle {
				return true
			}

			so.l.Lock()
			if so.b == resIdle || so.b == resOccupy {
				so.l.Unlock()
				return true
			}

			if conn, ok := so.s.(*wsconn); ok {
				conn.close(nil)
				so.l.Unlock()
				conn.closewait()
				//! 这里应该没有问题
				return true
			}
			so.l.Unlock()

			return true
		})
		//-----------------------------------------------------------
		wsl.httpWait.Done()
	}()
	time.Sleep(time.Millisecond * 1)

	return nil
}

// Connect : 无效
func (wsl *wslisten) connect(addr string) error {
	return nil
}

// Close : 关闭
func (wsl *wslisten) close(lck *util.ReSpinLock) error {
	err := wsl.httpSrv.Close()
	return err
}

func (wsl *wslisten) closewait() {
	wsl.httpWait.Wait()
}

func (wsl *wslisten) httpListen(addr string) (net.Listener, error) {
	if addr == "" {
		addr = ":http"
	}

	return net.Listen("tcp", addr)
}
