#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include "_cgo_export.h"
#include <stdio.h>

#define GOLUA_PANIC_MSG_WARAPPER "golua_panic_msg_warapper"

static int tag = 0;
static const char *const hooknames[] = {"call", "return", "line", "count", "tail return"};
static int hook_index = -1;
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


void hook(lua_State *L, lua_Debug *ar)
{
	int event;

	lua_pushlightuserdata(L, &hook_index);
	lua_rawget(L, LUA_REGISTRYINDEX);

	event = ar->event;
	lua_pushstring(L, hooknames[event]);

	lua_getinfo(L, "nG", ar);
	if (*(ar->what) == 'O') {
		lua_pushfstring(L, "[?%s]", ar->name);
	} else {
		lua_pushfstring(L, "%s:%d", ar->short_src, ar->linedefined > 0 ? ar->linedefined : 0);
	}

	lua_call(L, 2, 0);
}

lua_State* mlua_newstate(void* goallocf) {
  return lua_newstate(&go_alloc_wrapper, goallocf);
}

void mlua_setgostate(lua_State* L, size_t gohandle)
{
	lua_pushlightuserdata(L,(void*)&go_state_registry_key);
	lua_pushlightuserdata(L, (void*)gohandle);
	lua_settable(L, LUA_REGISTRYINDEX);
}

size_t mlua_getgostate(lua_State* L)
{
	size_t gohandle;
	//get gostate from registry entry
	lua_pushlightuserdata(L,(void*)&go_state_registry_key);
	lua_gettable(L, LUA_REGISTRYINDEX);
	gohandle = (size_t)lua_touserdata(L,-1);
	lua_pop(L,1);
	return gohandle;
}


int mlua_loadfile(lua_State *L, const char *filename) {
	return luaL_loadfile(L, filename);
}

mlua_loadbuffer(lua_State *L, const char *buffer, size_t sz, const char* name){
	return luaL_loadbuffer(L, buffer, sz, name);
}


static void call_ret_hook(lua_State *L) {
	lua_Debug ar;

	if (lua_gethook(L)) {
		lua_getstack(L, 0, &ar);
		lua_getinfo(L, "n", &ar);

		lua_pushlightuserdata(L, &hook_index);
		lua_rawget(L, LUA_REGISTRYINDEX);

		if (lua_type(L, -1) != LUA_TFUNCTION){
			lua_pop(L, 1);
			return;
        }

		lua_pushliteral(L, "return");
		lua_pushfstring(L, "[?%s]", ar.name);
		lua_pushliteral(L, "[GO]");

		lua_sethook(L, 0, 0, 0);
		lua_call(L, 3, 0);
		lua_sethook(L, hook, LUA_MASKCALL | LUA_MASKRET, 0);
	}
}

static int profiler_set_hook(lua_State *L) {
	if (lua_isnoneornil(L, 1)) {
		lua_pushlightuserdata(L, &hook_index);
		lua_pushnil(L);
		lua_rawset(L, LUA_REGISTRYINDEX);

		lua_sethook(L, 0, 0, 0);
	} else {
		luaL_checktype(L, 1, LUA_TFUNCTION);
		lua_pushlightuserdata(L, &hook_index);
		lua_pushvalue(L, 1);
		lua_rawset(L, LUA_REGISTRYINDEX);
		lua_sethook(L, hook, LUA_MASKCALL | LUA_MASKRET, 0);
	}
	return 0;
}

static int go_function_wrapper_wrapper(lua_State *L) {
	int ret;
	size_t gohandle = mlua_getgostate(L);

	ret = golua_call_gofunction(gohandle, 	mlua_tointeger(L, lua_upvalueindex(1)));

	if (lua_toboolean(L, lua_upvalueindex(2)))
  {
      lua_pushboolean(L, 0);
      lua_replace(L, lua_upvalueindex(2));
      return lua_error(L);
  }

	if (lua_gethook(L)) {
		call_ret_hook(L);
	}
}

void mlua_push_go_wrapper(lua_State* L, unsigned int wrapperid){
	lua_pushinteger(L, wrapperid);
	lua_pushboolean(L, 0);
	lua_pushcclosure(L, go_function_wrapper_wrapper, 2);
}

int panic_msg_warapper(lua_State *L){
	size_t gostateindex = mlua_getgostate(L);
  golua_panic_msg_func(gostateindex, (char*)lua_tolstring(L, -, NULL));
	return 0;
}

int mlua_pcall(lua_State* L, int nargs, int nresults, int errfunc){
	return lua_pcallk(L, nargs, nresults, errfunc, 0, NULL);
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

static const luaL_Reg mlualib[] = {
	{"sethook", profiler_set_hook},
	{NULL, NULL}
};

void luaopen_mlua(lua_State *L) {
  luaL_openlibs(L);
	luaL_newlib(L, mlualib);
	lua_setglobal(L, "mlua");
	lua_register(L, GOLUA_PANIC_MSG_WARAPPER, &panic_msg_warapper);
}
