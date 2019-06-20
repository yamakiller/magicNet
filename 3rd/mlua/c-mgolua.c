#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include "_cgo_export.h"
#include <stdio.h>

#define GOLUA_PANIC_MSG_WARAPPER "golua_panicmsg_warapper"
#define GOLUA_STATE_SELF "golua_state_self"


static const char go_state_registry_key = 'k';

typedef struct {
	unsigned int fake_id;
} GoStruct;

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

void mlua_setgostate(lua_State* L, void *goluaState)
{
	lua_pushlightuserdata(L,(void*)&go_state_registry_key);
	lua_pushlightuserdata(L, goluaState);
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

int mlua_getmetatable(lua_State *L, const char *k) {
	return lua_getfield(L, LUA_REGISTRYINDEX, k);
}

static int go_function_wrapper_wrapper(lua_State *L) {
	int ret;
	ret = golua_call_gofunction(mlua_getgostate(L), 	mlua_tointeger(L, lua_upvalueindex(1)));

	if (lua_toboolean(L, lua_upvalueindex(2)))
  {
      lua_pushboolean(L, 0);
      lua_replace(L, lua_upvalueindex(2));
      return lua_error(L);
  }

	return ret;
}

void mlua_push_go_wrapper(lua_State* L, unsigned int wrapperid) {
	lua_pushinteger(L, wrapperid);
	lua_pushboolean(L, 0);
	lua_pushcclosure(L, go_function_wrapper_wrapper, 2);
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

mlua_pushgostruct(lua_State *L, unsigned int wrapperid) {
	unsigned int* iidptr = (unsigned int *)lua_newuserdata(L, sizeof(unsigned int));
	*iidptr = wrapperid;
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
