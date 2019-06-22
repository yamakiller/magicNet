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
	"fmt"
	"unsafe"
	"reflect"
	"bytes"
	"encoding/gob"
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
	_name         string
	_source       string
	_short_source string
	_current_line int
}

func newState(L *C.lua_State) *State {
	newstate := &State{L, make([]interface{}, 0, 8), make([]uint, 0, 8)}
	C.mlua_setgostate(L, C.uintptr_t(uintptr(unsafe.Pointer(newstate))))
	return newstate
}

func (L *State) RegisterGoStruct(d interface{}) {
	gob.Register(d)
}

// lua_absindex
func (L *State) AbsIndex(index int) int {
	return int(C.lua_absindex(L._s, C.int(index)))
}

// lua_copy
func (L *State) Copy(fromindex int, toindex int) {
	C.lua_copy(L._s, C.int(fromindex), C.int(toindex))
}

// lua_checkstack
func (L *State) CheckStack(n int) int {
	return int(C.lua_checkstack(L._s, C.int(n)))
}

// lua_type
func (L *State) Type(index int) int {
	return int(C.lua_type(L._s, C.int(index)))
}

// lua_typename
func (L *State) TypeName(tp int) string {
	return C.GoString(C.lua_typename(L._s, C.int(tp)))
}

// lua_gettop
func (L *State) GetTop() int {
	return int(C.lua_gettop(L._s))
}

// lua_settop
func (L *State) SetTop(index int) {
	C.lua_settop(L._s, C.int(index))
}

// lua_pop
func (L *State) Pop(n int) {
	C.lua_settop(L._s, C.int(-n-1))
}

// lua_insert
func (L *State) Insert(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(1))
}

// lua_remove
func (L *State) Remove(index int) {
	C.lua_rotate(L._s, C.int(index), C.int(-1))
	L.Pop(1)
}

// lua_replace
func (L *State) Replace(index int) {
	C.mlua_replace(L._s, C.int(index))
}

// lua_pushboolean
func (L *State) PushBoolean(b bool) {
	var bint int
	if b {
		bint = 1
	} else {
		bint = 0
	}
	C.lua_pushboolean(L._s, C.int(bint))
}

// lua_pushstring
func (L *State) PushString(str string) {
	Cstr := C.CString(str)
	defer C.free(unsafe.Pointer(Cstr))
	C.lua_pushlstring(L._s, Cstr, C.size_t(len(str)))
}

func (L *State) PushBytes(b []byte) {
	C.lua_pushlstring(L._s, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
}

// lua_pushinteger
func (L *State) PushInteger(n int64) {
	C.lua_pushinteger(L._s, C.lua_Integer(n))
}

// lua_pushnil
func (L *State) PushNil() {
	C.lua_pushnil(L._s)
}

// lua_pushnumber
func (L *State) PushNumber(n float64) {
	C.lua_pushnumber(L._s, C.lua_Number(n))
}

// lua_pushthread
func (L *State) PushThread() (isMain bool) {
	return C.lua_pushthread(L._s) != 0
}

// lua_pushvalue
func (L *State) PushValue(index int) {
	C.lua_pushvalue(L._s, C.int(index))
}

// lua_pushcfunction -> PushGoFunction
func (L *State) PushGoFunction(f LuaGoFunction) {
	C.mlua_push_go_wrapper(L._s, unsafe.Pointer(&f))
}

// lua_pushcclosure -> PushGoClosure
func (L *State) PushGoClosure(f LuaGoFunction, args ...interface{}) {
  var argsNum int = 1
	C.lua_pushlightuserdata(L._s, unsafe.Pointer(&f))
	for _, val := range args {
		argsNum += 1
		switch(reflect.TypeOf(val).Kind()) {
		case reflect.Uint64:
		case reflect.Uint32:
		case reflect.Uint:
		case reflect.Int64:
		case reflect.Int32:
		case reflect.Int:
			L.PushInteger(reflect.ValueOf(val).Int())
			break;
		case reflect.Float64:
		case reflect.Float32:
			L.PushNumber(reflect.ValueOf(val).Float())
			break;
		case reflect.String:
			L.PushString(reflect.ValueOf(val).String())
			break;
		case reflect.Struct:
			L.PushUserGoStruct(val)
			break;
		case reflect.Uintptr:
		case reflect.UnsafePointer:
			L.PushLightGoStruct(unsafe.Pointer(reflect.ValueOf(val).Pointer()))
		  break;
		case reflect.Bool:
			L.PushBoolean(reflect.ValueOf(val).Bool())
			break;
		default:
			panic(fmt.Sprintf("mlua go Closure %s Type not supported", reflect.TypeOf(val).Name()))
			break;
		}
	}
	C.mlua_push_go_closure_wrapper(L._s, C.int(argsNum))
}

//----------------------------------------------------//
// mlua_upvalueindex
// 确保闭包参数从第二个位置开始防蚊
//
//----------------------------------------------------//
func (L *State) UpvalueIndex(n int) int {
	return int(C.mlua_upvalueindex(C.int(n)))
}

//----------------------------------------------------//
// mlua_pushgostruct => lua_newuserdata               //
// 内存管理权交由 lua虚拟机管理                           //
// 内存消耗略大                                         //
//----------------------------------------------------//
func (L *State) PushUserGoStruct(d interface{}) {
	var dby bytes.Buffer
	enc := gob.NewEncoder(&dby)
	err := enc.Encode(d)
	if (err != nil){
		panic(err)
		return
	}

	C.mlua_pushugostruct(L._s, (*C.char)(unsafe.Pointer(&dby.Bytes()[0])), C.size_t(len(dby.Bytes())))
}

// lua_pushlightuserdata =>nmlua_pushlgostruct
func (L *State) PushLightGoStruct(d unsafe.Pointer) {
	C.mlua_pushlgostruct(L._s, C.uintptr_t(uintptr(d)))
}

// lua_setglobal
func (L *State) SetGlobal(name string) {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	C.lua_setglobal(L._s, Cname)
}

// lua_getglobal
func (L *State) GetGlobal(name string) {
	Cname := C.CString(name)
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
	b := C.lua_tolstring(L._s, C.int(index), &size)
	return C.GoBytes(unsafe.Pointer(b), C.int(size))
}

// lua_tointeger
func (L *State) ToInteger(index int) int {
	return int(C.mlua_tointeger(L._s, C.int(index)))
}

// lua_tonumber
func (L *State) ToNumber(index int) float64 {
	return float64(C.mlua_tonumber(L._s, C.int(index)))
}

//----------------------------------------------
// lua_tougostruct => lua_touserdata
// 获取索引中的Go Struct 结构
// TODO： 思考感觉性能消耗不小!  没有没改进方案呢？
//----------------------------------------------
func (L *State) ToUserGoStruct(index int, s interface{}){
	r := (*C.struct_GoStruct)(C.mlua_tougostruct(L._s, C.int(index)))
	n := int(r._sz)
	d := bytes.NewBuffer(C.GoBytes(unsafe.Pointer(&r._data[0]), C.int(n)))
	dec := gob.NewDecoder(d)
	err := dec.Decode(s)
	if (err != nil) {
		panic(err)
	}
}

func (L *State) ToLightGoStruct(index int) unsafe.Pointer {
	return unsafe.Pointer(C.mlua_tolgostruct(L._s, C.int(index)))
}

// lua_rawlen
func (L *State) RawLen(index int) uint {
	return uint(C.lua_rawlen(L._s, C.int(index)))
}

// lua_topointer
func (L *State) ToPointer(index int) unsafe.Pointer {
	return unsafe.Pointer(C.lua_topointer(L._s, C.int(index)))
}

// lua_rawequal
func (L *State) RawEqual(index1 int, index2 int) int {
	return int(C.lua_rawequal(L._s, C.int(index1), C.int(index2)))
}

// lua_gettable
func (L *State) GetTable(index int) int {
	return int(C.lua_gettable(L._s, C.int(index)))
}

// lua_getfield
func (L *State) GetField(index int, k string) int {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	return int(C.lua_getfield(L._s, C.int(index), Ck))
}

// lua_geti
func (L *State) GetI(index int, n int64) int {
	return int(C.lua_geti(L._s, C.int(index), C.lua_Integer(n)))
}

// lua_rawget
func (L *State) RawGet(index int) int {
	return int(C.lua_rawget(L._s, C.int(index)))
}

// lua_rawgeti
func (L *State) RawGetI(index int, n int64) int {
	return int(C.lua_rawgeti(L._s, C.int(index), C.lua_Integer(n)))
}

// lua_rawgetp
func (L *State) RawGetP(index int, p unsafe.Pointer) int {
	return int(C.lua_rawgetp(L._s, C.int(index), p))
}

// lua_getmetatable
func (L *State) GetMetaTable(objindex int) int {
	return int(C.lua_getmetatable(L._s, C.int(objindex)))
}

// lua_getuservalue
func (L *State) GetUserValue(index int) int {
	return int(C.lua_getuservalue(L._s, C.int(index)))
}

// lua_settable
func (L *State) SetTable(index int) {
	C.lua_settable(L._s, C.int(index))
}

// lua_setfield
func (L *State) SetField(index int, k string) {
	Ck := C.CString(k)
	defer C.free(unsafe.Pointer(Ck))
	C.lua_setfield(L._s, C.int(index), Ck)
}

// lua_seti
func (L *State) SetI(index int, n int64) {
	C.lua_seti(L._s, C.int(index), C.lua_Integer(n))
}

// lua_rawset
func (L *State) RawSet(index int) {
	C.lua_rawset(L._s, C.int(index))
}

// lua_rawseti
func (L *State) RawSetI(index int, n int64) {
	C.lua_rawseti(L._s, C.int(index), C.lua_Integer(n))
}

// lua_rawsetp
func (L *State) RawSetP(index int, p unsafe.Pointer) {
	C.lua_rawsetp(L._s, C.int(index), p)
}

// lua_setmetatable
func (L *State) SetMetaTable(objindex int) {
	C.lua_setmetatable(L._s, C.int(objindex))
}

//lua_setuservalue
func (L *State) SetUserValue(index int) {
	C.lua_setuservalue(L._s, C.int(index))
}

// lua_gc
func (L *State) Gc(what int, data int) int {
	return int(C.lua_gc(L._s, C.int(what), C.int(data)))
}

// luaL_error
func (L *State) Error(sfmt string, v ...interface{}) int {
	Cerror := C.CString(fmt.Sprintf(sfmt, v...))
	defer C.free(unsafe.Pointer(Cerror))
	return int(C.mlua_error(L._s, Cerror))
}

// lua_concat
func (L *State) Concat(n int) {
	C.lua_concat(L._s, C.int(n))
}

// lua_setallocf
func (L *State) SetAllocF(f Alloc) {
	C.mlua_setallocf(L._s, unsafe.Pointer(&f))
}

// lua_pushglobaltable
func (L *State) PushGlobalTable() {
	C.mlua_pushglobaltable(L._s)
}

// lua_pcall
func (L *State) pcall(nargs, nresults, errfunc int) int {
	return int(C.mlua_pcall(L._s, C.int(nargs), C.int(nresults), C.int(errfunc)))
}

func (L *State) call_ex(nargs int, nresults int, catch bool) (err error) {
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
	C.lua_close(L._s)
}

func (L *State) StackTrace() []LuaStackEntry {
	r := []LuaStackEntry{}
	var d C.lua_Debug
	Sln := C.CString("Sln")
	defer C.free(unsafe.Pointer(Sln))

	for depth := 0; C.lua_getstack(L._s, C.int(depth), &d) > 0; depth++ {
		C.lua_getinfo(L._s, Sln, &d)
		ssb := make([]byte, C.LUA_IDSIZE)
		for i := 0; i < C.LUA_IDSIZE; i++ {
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
