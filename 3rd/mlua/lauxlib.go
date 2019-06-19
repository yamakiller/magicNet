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
	_code        int
	_message     string
	_stack_trace []LuaStackEntry
}

type LuaBuffer = C.luaL_Buffer

func (err *LuaError) Error() string {
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

// Executes file, returns nil for no errors or the lua error string on failure
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
func (L *State) DoString(s string) error {
	if r := L.LoadString(s); r != 0 {
		return &LuaError{r, L.ToString(-1), L.StackTrace()}
	}
	return L.Call(0, LUA_MULTRET)
}

// luaL_loadbuffer
func (L *State) LoadBuffer(data *byte, sz uint, name string) int {
	Cname := C.CString(name)
	defer C.free(unsafe.Pointer(Cname))
	return int(C.mlua_loadbuffer(L._s, (*C.char)((unsafe.Pointer)(data)), C.size_t(sz), Cname))
}

// luaL_argcheck
// WARNING: before b30b2c62c6712c6683a9d22ff0abfa54c8267863 the function ArgCheck had the opposite behaviour
func (L *State) Argcheck(cond bool, narg int, extramsg string) {
	if !cond {
		Cextramsg := C.CString(extramsg)
		defer C.free(unsafe.Pointer(Cextramsg))
		C.luaL_argerror(L._s, C.int(narg), Cextramsg)
	}
}

// luaL_argerror
func (L *State) ArgError(narg int, extramsg string) int {
	Cextramsg := C.CString(extramsg)
	defer C.free(unsafe.Pointer(Cextramsg))
	return int(C.luaL_argerror(L._s, C.int(narg), Cextramsg))
}

// luaL_callmeta
func (L *State) CallMeta(obj int, e string) int {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return int(C.luaL_callmeta(L._s, C.int(obj), Ce))
}

// Returns true if the value at index is light user data
func (L *State) IsLightUserdata(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUA_TLIGHTUSERDATA
}

// lua_isnil
func (L *State) IsNil(index int) bool { return LuaValType(C.lua_type(L._s, C.int(index))) == LUA_TNIL }

// lua_isnone
func (L *State) IsNone(index int) bool { return LuaValType(C.lua_type(L._s, C.int(index))) == LUA_TNONE }

// lua_isnoneornil
func (L *State) IsNoneOrNil(index int) bool { return int(C.lua_type(L._s, C.int(index))) <= 0 }

// lua_isnumber
func (L *State) IsNumber(index int) bool { return C.lua_isnumber(L._s, C.int(index)) == 1 }

// lua_isstring
func (L *State) IsString(index int) bool { return C.lua_isstring(L._s, C.int(index)) == 1 }

// lua_iscfunction
func (L *State) IsCFunction(index int) bool { return C.lua_iscfunction(L._s, C.int(index)) == 1 }

// lua_istable
func (L *State) IsTable(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUA_TTABLE
}

// lua_isthread
func (L *State) IsThread(index int) bool {
	return LuaValType(C.lua_type(L._s, C.int(index))) == LUA_TTHREAD
}

// lua_isuserdata
func (L *State) IsUserdata(index int) bool { return C.lua_isuserdata(L._s, C.int(index)) == 1 }

// lua_newtable
func (L *State) NewTable() {
	C.lua_createtable(L._s, 0, 0)
}

// lua_newuserdata
func (L *State) NewUserData(sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.lua_newuserdata(L._s, C.size_t(sz)))
}

// lua_newthread
func (L *State) NewThread() *State { //TODO: should have same lists as parent
	//		but may complicate gc
	s := C.lua_newthread(L._s)
	return &State{s, nil, nil}
}

// lua_next
func (L *State) Next(index int) int {
	return int(C.lua_next(L._s, C.int(index)))
}

// luaL_checkany
func (L *State) CheckAny(narg int) {
	C.luaL_checkany(L._s, C.int(narg))
}

// luaL_checkinteger
func (L *State) CheckInteger(narg int) int {
	return int(C.luaL_checkinteger(L._s, C.int(narg)))
}

// luaL_checknumber
func (L *State) CheckNumber(narg int) float64 {
	return float64(C.luaL_checknumber(L._s, C.int(narg)))
}

// luaL_checkstring
func (L *State) CheckString(narg int) string {
	var length C.size_t
	return C.GoString(C.luaL_checklstring(L._s, C.int(narg), &length))
}

// luaL_checktype
func (L *State) CheckType(narg int, t LuaValType) {
	C.luaL_checktype(L._s, C.int(narg), C.int(t))
}

// luaL_testudata
func (L *State) TestUData(narg int, tname string) unsafe.Pointer {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_testudata(L._s, C.int(narg), Ctname))
}

// luaL_checkudata
func (L *State) CheckUdata(narg int, tname string) unsafe.Pointer {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return unsafe.Pointer(C.luaL_checkudata(L._s, C.int(narg), Ctname))
}

// luaL_len
func (L *State) Len(index int) int {
	return int(C.luaL_len(L._s, C.int(index)))
}

// luaL_gsub
func (L *State) GSub(s string, p string, r string) string {
	Cs := C.CString(s)
	Cp := C.CString(p)
	Cr := C.CString(r)

	defer func() {
		C.free(unsafe.Pointer(Cs))
		C.free(unsafe.Pointer(Cp))
		C.free(unsafe.Pointer(Cr))
	}()

	return C.GoString(C.luaL_gsub(L._s, Cs, Cp, Cr))
}

// luaL_getsubtable
func (L *State) GetSubTable(index int, fname string) int {
	Cfname := C.CString(fname)
	defer C.free(unsafe.Pointer(Cfname))
	return int(C.luaL_getsubtable(L._s, C.int(index), Cfname))
}

// luaL_traceback
func (L *State) TraceBack(L1 *State, msg string, level int) {
	Cmsg := C.CString(msg)
	defer C.free(unsafe.Pointer(Cmsg))
	C.luaL_traceback(L._s, L1._s, Cmsg, C.int(level))
}

// luaL_getmetafield
func (L *State) GetMetaField(obj int, e string) bool {
	Ce := C.CString(e)
	defer C.free(unsafe.Pointer(Ce))
	return C.luaL_getmetafield(L._s, C.int(obj), Ce) != 0
}

// luaL_newmetatable
func (L *State) NewMetaTable(tname string) bool {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return C.luaL_newmetatable(L._s, Ctname) != 0
}

// luaL_setmetatable
func (L *State) SetMetatable(tname string) {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	C.luaL_setmetatable(L._s, Ctname)
}

func (L *State) GetMetatable(tname string) int {
	Ctname := C.CString(tname)
	defer C.free(unsafe.Pointer(Ctname))
	return int(C.mlua_getmetatable(L._s, Ctname))
}

// luaL_optinteger
func (L *State) OptInteger(narg int, d int) int {
	return int(C.luaL_optinteger(L._s, C.int(narg), C.lua_Integer(d)))
}

// luaL_optnumber
func (L *State) OptNumber(narg int, d float64) float64 {
	return float64(C.luaL_optnumber(L._s, C.int(narg), C.lua_Number(d)))
}

// luaL_optstring
func (L *State) OptString(narg int, d string) string {
	var length C.size_t
	Cd := C.CString(d)
	defer C.free(unsafe.Pointer(Cd))
	return C.GoString(C.luaL_optlstring(L._s, C.int(narg), Cd, &length))
}

// luaL_ref
func (L *State) Ref(t int) int {
	return int(C.luaL_ref(L._s, C.int(t)))
}

// luaL_typename
func (L *State) LTypename(index int) string {
	return C.GoString(C.lua_typename(L._s, C.lua_type(L._s, C.int(index))))
}

// luaL_unref
func (L *State) Unref(t int, ref int) {
	C.luaL_unref(L._s, C.int(t), C.int(ref))
}

// luaL_where
func (L *State) Where(lvl int) {
	C.luaL_where(L._s, C.int(lvl))
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

// luaL_buffinit
func (L *State) BuffInit(b *LuaBuffer) {
	C.luaL_buffinit(L._s, b)
}

// luaL_addlstring
func (L *State) AddLString(b *LuaBuffer, s unsafe.Pointer, sz uint) {
	C.luaL_addlstring(b, (*C.char)(s), C.size_t(sz))
}

// luaL_addstring
func (L *State) AddString(b *LuaBuffer, s string) {
	Cs := C.CString(s)
	defer C.free(unsafe.Pointer(Cs))

	C.luaL_addstring(b, Cs)
}

// luaL_addvalue
func (L *State) AddValue(b *LuaBuffer) {
	C.luaL_addvalue(b)
}

// luaL_pushresult
func (L *State) PushResult(b *LuaBuffer) {
	C.luaL_pushresult(b)
}

// luaL_pushresultsize
func (L *State) PushResultSize(b *LuaBuffer, sz uint) {
	C.luaL_pushresultsize(b, C.size_t(sz))
}

// luaL_buffinitsize
func (L *State) BuffInitSize(b *LuaBuffer, sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_buffinitsize(L._s, b, C.size_t(sz)))
}

// luaL_prepbuffsize
func (L *State) PrepBuffSize(b *LuaBuffer, sz uint) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_prepbuffsize(b, C.size_t(sz)))
}

// luaL_prepbuffer
func (L *State) PrepBuffer(b *LuaBuffer) unsafe.Pointer {
	return unsafe.Pointer(C.luaL_prepbuffsize(b, C.size_t(C.mlua_buffersize())))
}
