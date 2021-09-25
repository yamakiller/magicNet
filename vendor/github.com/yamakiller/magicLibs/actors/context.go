package actors

import (
	"fmt"
	"runtime"

	"github.com/yamakiller/magicLibs/actors/messages"
)

type contextState int32

const (
	stateNone contextState = iota
	stateAlive
	stateRestarting
	stateStopping
	stateStopped
)

func spawnContext(parent *Core, actor Actor, pid *PID) *Context {
	return &Context{_actor: actor,
		_pid:    pid,
		_parent: parent,
		_state:  stateAlive}
}

//Context actor 上下文
type Context struct {
	_actor   Actor
	_pid     *PID
	_parent  *Core
	_message interface{}
	_state   contextState
}

//Self 返回Actor引用对象
func (slf *Context) Self() *PID {
	return slf._pid
}

//Actor 返回Actor对象
func (slf *Context) Actor() Actor {
	return slf._actor
}

//Sender 返回当前消息的发送者
func (slf *Context) Sender() *PID {
	return UnWrapPackSender(slf._message)
}

//Message 返回当前消息
func (slf *Context) Message() interface{} {
	return UnWrapPackMessage(slf._message)
}

//MessageHeader 返回消息头
func (slf *Context) MessageHeader() ReadOnlyMessageHeader {
	return UnWrapPackHeader(slf._message)
}

//Info ...
func (slf *Context) Info(sfmt string, args ...interface{}) {
	slf._parent._log.Info(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Debug ...
func (slf *Context) Debug(sfmt string, args ...interface{}) {
	slf._parent._log.Debug(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Error ...
func (slf *Context) Error(sfmt string, args ...interface{}) {
	slf._parent._log.Error(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Warning ...
func (slf *Context) Warning(sfmt string, args ...interface{}) {
	slf._parent._log.Warning(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Fatal ...
func (slf *Context) Fatal(sfmt string, args ...interface{}) {
	slf._parent._log.Fatal(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Panic ...
func (slf *Context) Panic(sfmt string, args ...interface{}) {
	slf._parent._log.Panic(fmt.Sprintf("[%s]", slf._pid.ToString()), sfmt, args...)
}

//Post 提交用户消息
func (slf *Context) Post(pid *PID, message interface{}) {
	slf.postUsrMessage(pid, message)
}

//Request 请求消息
func (slf *Context) Request(pid *PID, message interface{}) {
	e := &Pack{
		Header:  nil,
		Message: message,
		Sender:  slf.Self(),
	}

	slf.postUsrMessage(pid, e)
}

//Stop 停止ACTOR
func (slf *Context) Stop(pid *PID) {
	pid.ref().Stop(pid)
}

func (slf *Context) postUsrMessage(pid *PID, message interface{}) {
	pid.postUsrMessage(message)
}

func (slf *Context) invokeSysMessage(message interface{}) {
	switch msg := message.(type) {
	case *messages.Started:
		slf.invokeUsrMessage(msg)
	case *messages.Stop:
		slf.onStop(msg)
	case *messages.Terminated:
		slf.onTerminated(msg)
	default:
		slf.Error("Message unfound:%+%v", msg)
	}
}

func (slf *Context) onStop(message interface{}) {
	if slf._state >= stateStopping {
		return
	}

	slf._state = stateStopping
	slf.invokeUsrMessage(messages.StoppingMessage)
	slf.tryTerminate()
}

func (slf *Context) onTerminated(message interface{}) {
	slf.invokeUsrMessage(message)
	slf.tryTerminate()
}

func (slf *Context) tryTerminate() {
	if slf._state == stateStopping {
		slf.finalizeStop()
	}
}

func (slf *Context) finalizeStop() {
	slf.invokeUsrMessage(messages.StoppedMessage)
	slf._state = stateStopped
	slf._parent.Delete(slf._pid)
}

func (slf *Context) invokeUsrMessage(message interface{}) {
	if slf._state == stateStopped {
		return
	}

	slf.processMessage(message)
}

func (slf *Context) processMessage(m interface{}) {
	slf._message = m
	slf.defaultReceive()
	slf._message = nil
}

func (slf *Context) receive(pack *Pack) {
	slf._message = pack
	slf.defaultReceive()
	slf._message = nil
}

func (slf *Context) defaultReceive() {
	if _, ok := slf.Message().(*messages.Kill); ok {
		slf.Stop(slf._pid)
		return
	}

	slf._actor.Receive(slf)
}

func (slf *Context) escalateFailure(reason interface{}, message interface{}) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	stackInfo := fmt.Sprintf("%s", buf[:n])

	slf.Panic("stack:\n %s\n", stackInfo)
	//slf.Panic("%+v, message:%+v", reason, message)
	slf.Self().postSysMessage(messages.SuspendMessage)
}
