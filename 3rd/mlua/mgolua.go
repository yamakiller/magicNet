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

type LuaGoAllocFunction func(ptr unsafe.Pointer, osize uint, nsize uint) unsafe.Pointer

type LuaGoHookFunction func(L *State, ar *LuaDebug)

type LuaGoFunction func(L *State) int

type State struct {
	_s *C.lua_State

	_h *LuaGoHookFunction
}

//export golua_call_allocf
func golua_call_allocf(fp uintptr, ptr uintptr, osize uint, nsize uint) uintptr {
	return uintptr((*((*LuaGoAllocFunction)(unsafe.Pointer(fp))))(unsafe.Pointer(ptr), osize, nsize))
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

//export golua_hook_gofunction
func golua_hook_gofunction(L unsafe.Pointer, ar *C.struct_lua_Debug) {
	L1 := (*State)(L)
  if L1._h == nil{
		return
	}

	var har LuaDebug
	har.Event = int(ar.event)
	har.Name  = C.GoString(ar.name)
	har.NameWhat = C.GoString(ar.namewhat)
	har.CurrentLine = int(ar.currentline)
	har.LineDefined = int(ar.linedefined)
	har.LastLineDefined = int(ar.lastlinedefined)
	har.NParams = uint8(ar.nparams)
	har.Nups = uint8(ar.nups)
	har.IsVararg = byte(ar.isvararg)
	har.IsTailCall = byte(ar.istailcall)
	har.ShortSrc = C.GoBytes(unsafe.Pointer(&ar.short_src[0]), C.int(LUA_IDSIZE))

	(*L1._h)(L1, &har)
}
