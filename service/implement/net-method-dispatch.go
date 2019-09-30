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
	Name   string
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
func (nmd *NetMethodDispatch) Register(key interface{}, f NetMethodFun) {
	nmd.RegisterType(reflect.TypeOf(key), f)
}

//RegisterType Registration network method, key reflect.Type
func (nmd *NetMethodDispatch) RegisterType(key reflect.Type, f NetMethodFun) {
	nmd.sync.Lock()
	defer nmd.sync.Unlock()
	nmd.m[key] = f
}

//Get Put back the network method according to the key object
func (nmd *NetMethodDispatch) Get(key interface{}) NetMethodFun {
	return nmd.GetType(reflect.TypeOf(key))
}

//GetType Put back the network method according to reflect.Type
func (nmd *NetMethodDispatch) GetType(key reflect.Type) NetMethodFun {
	nmd.sync.RLock()
	defer nmd.sync.RUnlock()
	f, success := nmd.m[key]
	if !success {
		return nil
	}
	return f
}

//Clear Clear all method maps
func (nmd *NetMethodDispatch) Clear() {
	nmd.sync.Lock()
	defer nmd.sync.Unlock()
	nmd.m = make(map[interface{}]NetMethodFun)
}
