#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include "_cgo_export.h"
#include <stdio.h>
#include <string.h>

#define GOLUA_PANIC_MSG_WARAPPER "golua_panicmsg_warapper"
#define GOLUA_STATE_SELF "golua_state_self"


static const char go_state_registry_key = 'k';
struct GoStruct{
	unsigned int _fakeId;
	size_t _sz;
  char _data[1];
};

int mlua_get_lib_version() {
	return 105;
}

void *go_alloc_wrapper(void *ud, void *ptr, size_t osize, size_t nsize){
  return (void*)golua_call_allocf((GoUintptr)ud, (GoUintptr)ptr, osize, nsize);
}

lua_State* mlua_newstate(void* goallocf) {
  return lua_newstate(&go_alloc_wrapper, goallocf);
}

void mlua_setallocf(lua_State* L, void* goallocf) {
	lua_setallocf(L, &go_alloc_wrapper, goallocf);
}

void mlua_setgostate(lua_State* L, uintptr_t goluaState)
{
	lua_pushlightuserdata(L,(void*)&go_state_registry_key);
	lua_pushlightuserdata(L, (void*)goluaState);
	lua_settable(L, LUA_REGISTRYINDEX);
}

void* mlua_getgostate(lua_State* L)
{
	void *goluaState;
	//get gostate from registry entry
	lua_pushlightuserdata(L,(void*)&go_state_registry_key);
	lua_gettable(L, LUA_REGISTRYINDEX);
	goluaState = lua_touserdata(L,-1);
	lua_pop(L,1);
	return goluaState;
}


int mlua_loadfile(lua_State *L, const char *filename) {
	return luaL_loadfile(L, filename);
}

int mlua_loadbuffer(lua_State *L, const char *buffer, size_t sz, const char* name){
	return luaL_loadbuffer(L, buffer, sz, name);
}

lua_Integer mlua_tointeger (lua_State *L, int idx) {
	return lua_tointeger(L, idx);
}

lua_Number mlua_tonumber(lua_State *L, int idx) {
	return lua_tonumber(L, idx);
}

const char *mlua_tostring(lua_State *L, int idx) {
	return lua_tostring(L, idx);
}

const void *mlua_tougostruct(lua_State *L, int idx) {
	struct GoStruct *pGoData = (struct GoStruct*)lua_touserdata(L, idx);
	return (const void *)pGoData;
}

const void *mlua_tolgostruct(lua_State *L, int idx) {
	return (const void *)lua_touserdata(L, idx);
}

int mlua_getmetatable(lua_State *L, const char *k) {
	return lua_getfield(L, LUA_REGISTRYINDEX, k);
}

static int go_function_wrapper_wrapper(lua_State *L) {
	return golua_call_gofunction(mlua_getgostate(L), 	(GoUintptr)lua_touserdata(L, lua_upvalueindex(1)));
}

int mlua_upvalueindex(int i) {
	return lua_upvalueindex(i + 1);
}

void mlua_push_fun_wrapper(lua_State* L, int nup) {
	lua_pushcclosure(L, go_function_wrapper_wrapper, nup);
}

void mlua_push_go_wrapper(lua_State* L,void* gofunc, int nup) {
	lua_pushlightuserdata(L, gofunc);
	int i;
	for(i = 0;i < nup;i++) {
		lua_pushvalue(L, -(nup + 1));
	}
	mlua_push_fun_wrapper(L, nup + 1);
	lua_pop(L, nup);
}

int panic_msg_warapper(lua_State *L) {
	golua_panicmsg_gofunction(mlua_getgostate(L), (char*)lua_tostring(L, -1));
	return 0;
}

int mlua_error(lua_State *L, const char *fmt) {
	return luaL_error(L, fmt);
}

void mlua_replace(lua_State *L, int idx) {
	lua_replace(L, idx);
}

void mlua_pushglobaltable(lua_State *L) {
	lua_pushglobaltable(L);
}

void mlua_pushugostruct(lua_State *L, char *godata, size_t sz) {
	struct GoStruct *pGoData = (struct GoStruct*)lua_newuserdata(L, sizeof(struct GoStruct) + sz);
	pGoData->_fakeId = 0;
	pGoData->_sz = sz;
	memcpy(&pGoData->_data[0], (void*)godata, sz);
}

void mlua_pushlgostruct(lua_State *L, uintptr_t p) {
	lua_pushlightuserdata(L, (void*)p);
}

unsigned int mlua_isgostruct(lua_State *L, int idx) {
	if (lua_isuserdata(L, idx) != 0) {
		unsigned int* iidptr = lua_touserdata(L, idx);
		return *iidptr;
	}
	return 0;
}

int mlua_pcall(lua_State* L, int nargs, int nresults, int errfunc){
	return lua_pcallk(L, nargs, nresults, errfunc, 0, NULL);
}


void luaopen_mlua(lua_State *L) {
  luaL_openlibs(L);
	lua_register(L, GOLUA_PANIC_MSG_WARAPPER, &panic_msg_warapper);
}

int mlua_buffersize() {
	return LUAL_BUFFERSIZE;
}

static void go_hook_wrapper(lua_State *L, lua_Debug *ar){
	golua_hook_gofunction(mlua_getgostate(L), ar);
}

void mlua_sethook(lua_State *L, int mask, int count) {
	lua_sethook(L, &go_hook_wrapper, mask, count);
}
