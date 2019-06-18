package mlua

/*
#cgo CFLAGS: -I ${SRCDIR}/lua

#include <lua.h>
#include <lualib.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

type Alloc func(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer

type LuaGoFunction func(L *State) int

type State struct {
	_s *C.lua_State

	_Index uintptr

	_registry []interface{}

	_freeIndices []uint
}

//export golua_call_allocf
func golua_call_allocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	return uintptr((*((*Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
}

//export golua_call_gofunction
func golua_call_gofunction(gostateindex uintptr, fid uint) int {
	L1 := getGoState(gostateindex)
	if fid < 0 {
		panic(&LuaError{0, "Requested execution of an unknown function", L1.StackTrace()})
	}
	f := L1._registry[fid].(LuaGoFunction)
	return f(L1)
}

//export golua_panicmsg_gofunction
func golua_panicmsg_gofunction(gostateindex uintptr, z *C.char) {
	L := getGoState(gostateindex)
	s := C.GoString(z)

	panic(&LuaError{LUA_ERRERR, s, L.StackTrace()})
}
