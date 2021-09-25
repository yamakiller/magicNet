package mkcp

/*
#cgo CFLAGS: -I ${SRCDIR}
#cgo lkcp LDFLAGS: -lkcp
#cgo kcpa LDFLAGS: -lkcp -lm -ldl
#cgo linux,!lkcp,!kcpa LDFLAGS: -L${SRCDIR} -lkcp
#cgo darwin,!lkcp,!kcpa LDFLAGS: -lkcp
#cgo freebsd,!kcpa LDFLAGS: -lkcp
#cgo windows,!lkcp LDFLAGS: -L${SRCDIR} -lkcp -lmingwex -lmingw32

#include <kcp/ikcp.h>
#include <stdlib.h>

#include "mgokcp.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

//KCPGoOutputFunc kcp go output function
type KCPGoOutputFunc func([]byte, interface{}) int32

//New 创建一个KCP
func New(conv uint32, user interface{}) *KCP {
	goKcp := &KCP{_usr: user}
	goKcp._kcp = C.mkcp_create(C.uint(conv), C.uintptr_t(uintptr(unsafe.Pointer(goKcp))))
	return goKcp
}

//Free 释放一个KCP
func Free(kcp *KCP) {
	C.mkcp_release(kcp._kcp)
}

//export go_output
func go_output(buf []byte, ptr uintptr) int {
	kcp := (*KCP)(unsafe.Pointer(ptr))
	return int(kcp._output(buf, kcp._usr))
}

//KCP Connection
type KCP struct {
	_kcp    *C.ikcpcb
	_usr    interface{}
	_output KCPGoOutputFunc
}

//WithOutput 设置输出回调函数
func (slf *KCP) WithOutput(output KCPGoOutputFunc) {
	slf._output = output
}

//User 返回KCPConn的表示参数
func (slf *KCP) User() interface{} {
	return slf._usr
}

//Recv 接收数据 buffer 保证内存连续性 make([]byte, n)
func (slf *KCP) Recv(buffer []byte, size int32) int32 {
	return int32(C.ikcp_recv(slf._kcp,
		(*C.char)(unsafe.Pointer(&buffer[0])),
		C.int(size)))
}

//Send 发送数据 buffer 保证内存连续性　make([]byte, n)
func (slf *KCP) Send(buffer []byte, size int32) (int32, error) {
	n := int32(C.mkcp_send(slf._kcp,
		(*C.char)(unsafe.Pointer(&buffer[0])),
		C.int(size)))
	if n == -1 {
		return 0, errors.New("Waiting to send data smaller than 0")
	} else if n == -2 {
		return 0, errors.New("Data window overflow")
	}

	return size, nil
}

//Update ....
func (slf *KCP) Update(current uint32) {
	C.ikcp_update(slf._kcp, C.uint(current))
}

//Check ...
func (slf *KCP) Check(current uint32) uint32 {
	return uint32(C.ikcp_check(slf._kcp, C.uint(current)))
}

//Input 对接收到的数据进行处理 保证data的连续性
func (slf *KCP) Input(data []byte, size int32) int32 {
	return int32(C.mkcp_input(slf._kcp,
		(*C.char)(unsafe.Pointer(&data[0])),
		C.int(size)))
}

//Flush 输出数据直接回调到KCP的output函数
func (slf *KCP) Flush() {
	C.ikcp_flush(slf._kcp)
}

//PeekSize ...
func (slf *KCP) PeekSize() int32 {
	return int32(C.ikcp_peeksize(slf._kcp))
}

//SetMTU default 1400
func (slf *KCP) SetMTU(mtu int32) int {
	return int(C.ikcp_setmtu(slf._kcp, C.int(mtu)))
}

//SetRxMinRto 设置rx_minrto
func (slf *KCP) SetRxMinRto(rxMinRto int32) {
	slf._kcp.rx_minrto = C.int(rxMinRto)
}

//SetFastResend 设置快速模式
func (slf *KCP) SetFastResend(fastResend int32) {
	slf._kcp.fastresend = C.int(fastResend)
}

//WndSize 色织窗口
func (slf *KCP) WndSize(sndWnd, rcvWnd int32) int32 {
	return int32(C.ikcp_wndsize(slf._kcp, C.int(sndWnd), C.int(rcvWnd)))
}

//WaitSnd ...
func (slf *KCP) WaitSnd() int32 {
	return int32(C.ikcp_waitsnd(slf._kcp))
}

//NoDelay 设置延迟
func (slf *KCP) NoDelay(nodelay, interval, resend, nc int32) int32 {
	return int32(C.ikcp_nodelay(slf._kcp,
		C.int(nodelay),
		C.int(interval),
		C.int(resend),
		C.int(nc)))
}

//GetConv 保证ptr的连续性 make([]byte, n)
func GetConv(ptr []byte) uint32 {
	return uint32(C.mkcp_getconv((*C.char)(unsafe.Pointer(&ptr[0]))))
}
