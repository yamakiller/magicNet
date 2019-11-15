package actor

import (
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/engine/logger"
)

/**********************************************************************************
* The upper 15 bits indicate the server address | Low 17 indicates the PID number *
***********************************************************************************/
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

//HexToID desc
//@method HexToID desc: Hexadecimal string, converted to uint32
//@param  (string) hex string
//@return (uint32)
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

//Address desc
//@method Address desc: Returns address information
//@param (uint32)
func (pid *PID) Address() uint32 {
	return pid.ID >> pidKeyBit
}

//Key desc
//@method Key desc: Returns unique ID Key without address information
//@param (uint32)
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
	return idToHex(pid.ID)
}

//Stop desc
//@method Stop desc: Stop PID and stop Actor
func (pid *PID) Stop() {
	pid.ref().Stop(pid)
}

//Tell desc
//@method Tell desc: transfer
//@param (interface{})
func (pid *PID) Tell(message interface{}) {
	ctx := DefaultSchedulerContext
	ctx.Send(pid, message)
}

//Request desc
//@method Request desc: Request a custom recovery target PID
//@param (interface{}) message
//@param (*PID)
func (pid *PID) Request(message interface{}, responseTo *PID) {
	ctx := DefaultSchedulerContext
	ctx.RequestWithCustomSender(pid, message, responseTo)
}

//RequestFuture desc
//@method RequestFuture desc: Request and wait for a reply [with timeout]
//@param (interface{}) message
//@param (time.Duration) time out
func (pid *PID) RequestFuture(message interface{}, timeOut time.Duration) *Future {
	ctx := DefaultSchedulerContext
	return ctx.RequestFuture(pid, message, timeOut)
}

//StopFuture desc
//@method StopFuture desc: Stop the PID and wait for a reply
//@return (*Future)
func (pid *PID) StopFuture() *Future {
	future := NewFuture(10 * time.Second)

	pid.sendSysMessage(&Watch{Watcher: future.pid})
	pid.Stop()
	return future
}

//StopWait desc
//@method StopWait desc: Stop and wait
func (pid *PID) StopWait() {
	pid.StopFuture().Wait()
}
