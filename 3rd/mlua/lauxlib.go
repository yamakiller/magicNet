package mlua

//#include <lua.h>
//#include <lauxlib.h>
//#include <lualib.h>
//#include <stdlib.h>
//#include "mgolua.h"
import "C"
import (
	"unsafe"
)

type LuaError struct {
	_code int
	_message string
	_stack_trace []LuaStackEntry
}

func (err *LuaError) GetWhat() string {
	return err._message
}

func (err *LuaError) GetCode() int {
	return err._code
}

func (err *LuaError) GetStackTrace() []LuaStackEntry {
	return err._stack_trace
}

// luaL_loadfile
func (L *State) LoadFile(filename string) int {
	Cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(Cfilename))
	return int(C.mlua_loadfile(L._s, Cfilename))
}

// luaL_dofile
func (L *State) DoFile(filename string) error {
	if r := L.LoadFile(filename); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUA_MULTRET)
}

// luaL_loadstring
func (L *State) LoadString(s string) int {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	return int(C.luaL_loadstring(L._s, Cs))
}

// luaL_dostring
func (L *State) DoString(s string) {
	if r := L.LoadString(s); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUA_MULTRET)
}

// luaL_loadbuffer
func (L *State) LoadBuffer(data *byte, uint sz, name string) int {
	Cname := C.CString(Cname)
	defer C.free(unsafe.Pointer(Cname))
	return int(C.mlua_loadbuffer(L._s, (*C.char)((unsafe.Pointer)(data)), C.size_t(sz), Cname))
}

// luaL_newstate
func NewState() *State {
	ls := (C.luaL_newstate())
	if ls == nil {
		return nil
	}
	L := newState(ls)
	return L
}

// luaL_openlibs
func (L *State) OpenLibs() {
	C.luaL_openlibs(L._s)
}
