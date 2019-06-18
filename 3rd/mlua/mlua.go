package mlua

/*
#cgo CFLAGS: -I ${SRCDIR}/lua
#cgo llua LDFLAGS: -llua
#cgo luaa LDFLAGS: -llua -lm -ldl
#cgo linux,!llua,!luaa LDFLAGS: -llua
#cgo darwin,!llua,!luaa LDFLAGS: -llua
#cgo freebsd,!luaa LDFLAGS: -llua
#cgo windows,!llua LDFLAGS: -L${SRCDIR} -llua -lmingwex -lmingw32

#include <lua.h>
#include <stdlib.h>

#include "mgolua.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

type LuaValType int

const (
	LUA_TNIL           = LuaValType(C.LUA_TNIL)
	LUA_TNUMBER        = LuaValType(C.LUA_TNUMBER)
	LUA_TBOOLEAN       = LuaValType(C.LUA_TBOOLEAN)
	LUA_TSTRING        = LuaValType(C.LUA_TSTRING)
	LUA_TTABLE         = LuaValType(C.LUA_TTABLE)
	LUA_TFUNCTION      = LuaValType(C.LUA_TFUNCTION)
	LUA_TUSERDATA      = LuaValType(C.LUA_TUSERDATA)
	LUA_TTHREAD        = LuaValType(C.LUA_TTHREAD)
	LUA_TLIGHTUSERDATA = LuaValType(C.LUA_TLIGHTUSERDATA)
)

const (
	LUA_VERSION       = C.LUA_VERSION
	LUA_RELEASE       = C.LUA_RELEASE
	LUA_VERSION_NUM   = C.LUA_VERSION_NUM
	LUA_COPYRIGHT     = C.LUA_COPYRIGHT
	LUA_AUTHORS       = C.LUA_AUTHORS
	LUA_MULTRET       = C.LUA_MULTRET
	LUA_REGISTRYINDEX = C.LUA_REGISTRYINDEX
	LUA_YIELD         = C.LUA_YIELD
	LUA_ERRRUN        = C.LUA_ERRRUN
	LUA_ERRSYNTAX     = C.LUA_ERRSYNTAX
	LUA_ERRMEM        = C.LUA_ERRMEM
	LUA_ERRERR        = C.LUA_ERRERR
	LUA_TNONE         = C.LUA_TNONE
	LUA_MINSTACK      = C.LUA_MINSTACK
	LUA_GCSTOP        = C.LUA_GCSTOP
	LUA_GCRESTART     = C.LUA_GCRESTART
	LUA_GCCOLLECT     = C.LUA_GCCOLLECT
	LUA_GCCOUNT       = C.LUA_GCCOUNT
	LUA_GCCOUNTB      = C.LUA_GCCOUNTB
	LUA_GCSTEP        = C.LUA_GCSTEP
	LUA_GCSETPAUSE    = C.LUA_GCSETPAUSE
	LUA_GCSETSTEPMUL  = C.LUA_GCSETSTEPMUL
	LUA_HOOKCALL      = C.LUA_HOOKCALL
	LUA_HOOKRET       = C.LUA_HOOKRET
	LUA_HOOKLINE      = C.LUA_HOOKLINE
	LUA_HOOKCOUNT     = C.LUA_HOOKCOUNT
	LUA_MASKCALL      = C.LUA_MASKCALL
	LUA_MASKRET       = C.LUA_MASKRET
	LUA_MASKLINE      = C.LUA_MASKLINE
	LUA_MASKCOUNT     = C.LUA_MASKCOUNT
	LUA_ERRFILE       = C.LUA_ERRFILE
	LUA_NOREF         = C.LUA_NOREF
	LUA_REFNIL        = C.LUA_REFNIL
	LUA_FILEHANDLE    = C.LUA_FILEHANDLE
	LUA_COLIBNAME     = C.LUA_COLIBNAME
	LUA_TABLIBNAME    = C.LUA_TABLIBNAME
	LUA_IOLIBNAME     = C.LUA_IOLIBNAME
	LUA_OSLIBNAME     = C.LUA_OSLIBNAME
	LUA_STRLIBNAME    = C.LUA_STRLIBNAME
	LUA_MATHLIBNAME   = C.LUA_MATHLIBNAME
	LUA_DBLIBNAME     = C.LUA_DBLIBNAME
	LUA_LOADLIBNAME   = C.LUA_LOADLIBNAME
)

type LuaStackEntry struct {
	_name string
	_source string
	_short_source string
	_current_line int
}

var goStates map[uintptr]*State
var goStatesMutex sync.Mutex

func init() {
	goStates = make(map[uintptr]*State, 16)
}

func registerGoState(L *State) {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	L._Index = uintptr(unsafe.Pointer(L))
	goStates[L._Index] = L
}

func unregisterGoState(L *State) {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	delete(goStates, L._Index)
}

func getGoState(gohandle uintptr) *State {
	goStatesMutex.Lock()
	defer goStatesMutex.Unlock()
	return goStates[gohandle]
}

func newState(L *C.lua_State) *State {
	newstate := &State{L, 0, make([]interface{}, 0, 8), make([]uint, 0, 8)}
	registerGoState(newstate)
	C.mlua_setgostate(L, C.size_t(newstate._Index))
	return newstate
}

func (L *State) addFreeIndex(i uint) {
	freelen := len(L._freeIndices)
	//reallocate if necessary
	if freelen+1 > cap(L._freeIndices) {
		newSlice := make([]uint, freelen, cap(L._freeIndices)*2)
		copy(newSlice, L._freeIndices)
		L._freeIndices = newSlice
	}
	//reslice
	L._freeIndices = L._freeIndices[0 : freelen+1]
	L._freeIndices[freelen] = i
}

func (L *State) getFreeIndex() (index uint, ok bool) {
	freelen := len(L._freeIndices)
	//if there exist entries in the freelist
	if freelen > 0 {
		i := L._freeIndices[freelen-1] //get index
		//fmt.Printf("Free indices before: %v\n", L.freeIndices)
		L._freeIndices = L._freeIndices[0 : freelen-1] //'pop' index from list
		//fmt.Printf("Free indices after: %v\n", L.freeIndices)
		return i, true
	}
	return 0, false
}

//returns the registered function id
func (L *State) register(f interface{}) uint {
	index, ok := L.getFreeIndex()
	if !ok {
		index = uint(len(L._registry))
		//reallocate backing array if necessary
		if index+1 > uint(cap(L._registry)) {
			newcap := cap(L._registry) * 2
			if index+1 > uint(newcap) {
				newcap = int(index + 1)
			}
			newSlice := make([]interface{}, index, newcap)
			copy(newSlice, L._registry)
			L._registry = newSlice
		}
		//reslice
		L._registry = L._registry[0 : index+1]
	}

	L._registry[index] = f
	return index
}

func (L *State) unregister(fid uint) {
	if (fid < uint(len(L._registry))) && (L._registry[fid] != nil) {
		L._registry[fid] = nil
		L.addFreeIndex(fid)
	}
}

func (L *State) PushGoFunction(f LuaGoFunction) {
	fid := L.register(f)
	C.mlua_push_go_wrapper(L._s, C.uint(fid))
}

// lua_gettop
func (L *State) GetTop() int {
		return int(C.lua_gettop(L._s))
}

// lua_insert
func (L *State) Insert(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(1))
}

// lua_remove
func (L *State) Remove(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(-1))
	C.lua_pop(L._s, C.int(1))
}

// lua_setglobal
func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setglobal(L._s, Cname)
}

// lua_getglobal
func（L *State） GetGlobal(name string) {
	Cname := C.CString(nae)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_getglobal(L._s, Cname)
}

// lua_tostring
func (L *State) ToString(index int) string {
	var size C.size_t
	r := C.lua_tolstring(L._s, C.int(index), &size)
	return C.GoStringN(r, C.int(size))
}

// luaL_tolstring
func (L *State) ToBytes(index int) []byte {
	var size C.size_t
	b := C.lua_tolstring(L._s, index, &size)
	return C.GoBytes(unsafe.Pointer(b), C.int(size))
}

// lua_tointeger
func (L *State) ToInteger(index int) int {
	return int(C.mlua_tointeger(L._s, index))
}

// lua_tonumber
func (L *State) ToNumber(index int) float64 {
	return float64(C.mlua_tonumber(L._s, C.int(index)))
}

// lua_pcall
func (L *State) pcall(nargs, nresults, errfunc int) {
	return int(C.mlua_pcall(L._s, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

func (L *State) call_ex(nargs int, nresults int, catch bool) (err error){
	if catch {
		defer func() {
			if err2 := recover(); err2 != nil {
				if _, ok := err2.(error); ok {
					err = err2.(error)
				}
				return
			}
		}()
	}

	L.GetGlobal(C.GOLUA_PANIC_MSG_WARAPPER)
	erridx := L.GetTop() - nargs - 1
	L.Insert(erridx)
	r := L.pcall(nargs, nresults, erridx)
	L.Remove(erridx)
	if r != 0 {
		err = &LuaError{r, L.ToString(-1), L.StackTrace()}
		if !catch {
			panic(err)
		}
	}
	return
}

// lua_call
func (L *State) Call(nargs, nresults int) (err error) {
	return L.call_ex(nargs, nresults, true)
}

// Registers a Go function as a global variable
func (L *State) Register(name string, f LuaGoFunction) {
	L.PushGoFunction(f)
	L.SetGlobal(name)
}

//exprot NewStateEnv
func NewStateAlloc(f Alloc) *State {
	ls := C.mlua_newstate(unsafe.Pointer(&f))
	return newState(ls)
}

// lua_close
func (L *State) Close() {
	unregisterGoState(L)
	C.lua_close(L._s)
}

func (L *State) StackTrace() []LuaStackEntry {
	r := []LuaStackEntry{}
	var d C.lua_Debug
	Sln := C.CString(Sln)
	defer C.free(unsafe.Pointer(Sln))

	for depth := 0; C.lua_getstack(L._s, C.int(depth), &d) > 0; depth++ {
		C.lua_getinfo(L._s, Sln, &d)
		ssb := make([]byte, C.LUA_IDSIZE)
		for i:= 0;i < C.LUA_IDSIZE; i++ {
			ssb[i] = byte(d.short_src[i])
			if ssb[i] == 0 {
				ssb = ssb[:i]
				break
			}
		}
		ss := string(ssb)
		r = append(r, LuaStackEntry{C.GoString(d.name), C.GoString(d.source), ss, int(d.currentline)})
	}
	return r
}

func (L *State) NewError(msg string) *LuaError {
	return &LuaError{0, msg, L.StackTrace()}
}
