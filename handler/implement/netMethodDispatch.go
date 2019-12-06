package implement

import (
	"reflect"
	"sync"
)

//INetMethodEvent Network Method Event interface
type INetMethodEvent interface {
}

//NetMethodEvent Network Method Event
type NetMethodEvent struct {
	Name   interface{}
	Socket int32
	Wrap   []byte
}

//NetMethodFun xxx
type NetMethodFun func(event INetMethodEvent)

//SpawnMethodDispatch Building a network method scheduler
func SpawnMethodDispatch() NetMethodDispatch {
	return NetMethodDispatch{m: make(map[interface{}]NetMethodFun)}
}

//NewMethodDispatch Newing a network method scheduler
func NewMethodDispatch() *NetMethodDispatch {
	return &NetMethodDispatch{m: make(map[interface{}]NetMethodFun)}
}

//NetMethodDispatch Network Method Scheduler
type NetMethodDispatch struct {
	m    map[interface{}]NetMethodFun
	sync sync.RWMutex
}

//Register Registration network method
func (slf *NetMethodDispatch) Register(key interface{}, f NetMethodFun) {
	slf.RegisterType(reflect.TypeOf(key), f)
}

//RegisterType Registration network method, key reflect.Type
func (slf *NetMethodDispatch) RegisterType(key reflect.Type, f NetMethodFun) {
	slf.sync.Lock()
	defer slf.sync.Unlock()
	slf.m[key] = f
}

//Get Put back the network method according to the key object
func (slf *NetMethodDispatch) Get(key interface{}) NetMethodFun {
	return slf.GetType(reflect.TypeOf(key))
}

//GetType Put back the network method according to reflect.Type
func (slf *NetMethodDispatch) GetType(key reflect.Type) NetMethodFun {
	slf.sync.RLock()
	defer slf.sync.RUnlock()
	f, success := slf.m[key]
	if !success {
		return nil
	}
	return f
}

//Clear Clear all method maps
func (slf *NetMethodDispatch) Clear() {
	slf.sync.Lock()
	defer slf.sync.Unlock()
	slf.m = make(map[interface{}]NetMethodFun)
}
