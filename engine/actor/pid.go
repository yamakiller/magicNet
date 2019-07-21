package actor

import (
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/util"
)

/***************************************
* 高15位表示服务器地址 | 低17表示PID编号 *
****************************************/
const (
	pidMask   = 0x1ffff
	pidMax    = pidMask
	pidKeyBit = 17
)

func idToHex(u uint32) string {
	const (
		digits = "0123456789ABCDEF"
	)

	var str [10]byte
	str[0] = '$'
	var i uint32
	for i = 0; i < 8; i++ {
		str[i+1] = digits[(u>>((7-i)*4))&0xf]
	}

	return string(str[:8])
}

// HexToID : 16进制字符串，转换为 uint32
func HexToID(hex string) uint32 {
	var i uint32
	var addr uint32
	var len = uint32(strings.Count(hex, "") - 1)
	for i = 1; i < len; i++ {
		c := hex[i]
		if c >= '0' && c <= '9' {
			c = c - '0'
		} else if c >= 'a' && c <= 'f' {
			c = c - 'a' + 10
		} else if c >= 'A' && c <= 'F' {
			c = c - 'A' + 10
		} else {
			util.Assert(false, "Id unknown character")
		}
		addr = addr*16 + uint32(c)
	}
	return addr
}

func pidFromID(id string, p *PID) {
	p.ID = HexToID(id)
}

func pidIsRemote(id uint32) bool {
	if (id >> pidKeyBit) == globalRegistry.GetLocalAddress() {
		return false
	}

	return true
}

// PID : Actor ID对象D
type PID struct {
	ID uint32
	p  *Process
}

// Address : 获取地址信息
func (pid *PID) Address() uint32 {
	return pid.ID >> pidKeyBit
}

// Key : 获取唯一的ID Key 不带地址信息
func (pid *PID) Key() uint32 {
	return pid.ID & pidMask
}

func (pid *PID) ref() Process {
	p := (*Process)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p))))
	if p != nil {
		if l, ok := (*p).(*AtrProcess); ok && atomic.LoadInt32(&l.death) == 1 {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), nil)
		} else {
			return *p
		}
	}

	ref, exits := globalRegistry.Get(pid)
	if exits {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&pid.p)), unsafe.Pointer(&ref))
	}

	return ref
}

func (pid *PID) sendUsrMessage(message interface{}) {
	ref := pid.ref()
	ref.SendUsrMessage(pid, message)
	overload := ref.OverloadUsrMessage()
	if overload > 0 {
		logger.Warning(pid.ID, "mailbox overload :%d", overload)
	}
}

func (pid *PID) sendSysMessage(message interface{}) {
	pid.ref().SendSysMessage(pid, message)
}

func (pid *PID) String() string {
	return ""
}

// Stop ： 停止PID既停止Actor
func (pid *PID) Stop() {
	pid.ref().Stop(pid)
}

// NewPID ：新建一个PID
/*func NewPID() *PID {
	pid := &PID{}
	globalRegistry.Register(pid)
	return pid
}*/

// Tell : 调用
func (pid *PID) Tell(message interface{}) {
	ctx := DefaultSchedulerContext
	ctx.Send(pid, message)
}

// Request ：请求自定义恢复目标PID
func (pid *PID) Request(message interface{}, responseTo *PID) {
	ctx := DefaultSchedulerContext
	ctx.RequestWithCustomSender(pid, message, responseTo)
}

// RequestFuture ：请求并等待回复[带超时]
func (pid *PID) RequestFuture(message interface{}, timeOut time.Duration) *Future {
	ctx := DefaultSchedulerContext
	return ctx.RequestFuture(pid, message, timeOut)
}

// StopFuture ：停止PID并等待回复
func (pid *PID) StopFuture() *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	pid.Stop()
	return future
}

// StopWait ：停止并等待
func (pid *PID) StopWait() {
	pid.StopFuture().Wait()
}
