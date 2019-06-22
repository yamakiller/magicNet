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

	_registry []interface{}

	_freeIndices []uint
}

//export golua_call_allocf
func golua_call_allocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	return uintptr((*((*Alloc)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
}

//export golua_call_gofunction
func golua_call_gofunction(L unsafe.Pointer, f uintptr) int {
	L1 := (*State)(L)
	return (*((*LuaGoFunction)(unsafe.Pointer(f))))(L1)
}

//export golua_panicmsg_gofunction
func golua_panicmsg_gofunction(L unsafe.Pointer, z *C.char) {
	L1 := (*State)(L)
	s := C.GoString(z)

	panic(&LuaError{LUA_ERRERR, s, L1.StackTrace()})
}
