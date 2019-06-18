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

// luaL_loadfile
func (L *State) LoadFile(filename string) int {
	Cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(Cfilename))
	return int(C.mlua_loadfile(L._s, Cfilename))
}

// luaL_loadstring
func (L *State) LoadString(s string) int {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))
	return int(C.luaL_loadstring(L._s, Cs))
}

// luaL_dostring
func (L *State) DoString(s string) {
	if L.LoadString(s) == 0 {
		C.mlua_pcall(L._s, C.int(0), C.int(LUA_MULTRET), C.int(0))
	}
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
