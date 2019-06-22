#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include <stdint.h>

#define GOLUA_PANIC_MSG_WARAPPER "golua_panicmsg_warapper"

struct GoStruct{
	unsigned int _fakeId;
	size_t _sz;
  char _data[1];
};

typedef int (*lua_GOWrapperCaller) (lua_State *L, unsigned int wrapperid, int top);

int mlua_get_lib_version();

lua_State* mlua_newstate(void* goallocf);

void mlua_setallocf(lua_State* L, void* goallocf);

void mlua_setgostate(lua_State *L, uintptr_t goluaState);

void* mlua_getgostate(lua_State* L);

int mlua_loadfile(lua_State *L, const char *filename);

int mlua_loadbuffer(lua_State *L, const char *buffer, size_t sz, const char* name);

void mlua_push_go_wrapper(lua_State* L, void* gofunc);

int mlua_pcall(lua_State* L, int nargs, int nresults, int errfunc);

lua_Integer mlua_tointeger(lua_State *L, int idx);

lua_Number mlua_tonumber(lua_State *L, int idx);

const char *mlua_tostring(lua_State *L, int idx);

const void *mlua_tougostruct(lua_State *L, int idx);

int mlua_error(lua_State *L, const char *fmt);

void mlua_replace(lua_State *L, int idx);

void mlua_pushglobaltable(lua_State *L);

void mlua_pushugostruct(lua_State *L, char *goStruct, size_t sz);

unsigned int mlua_isgostruct(lua_State *, int idx);

int mlua_getmetatable(lua_State *L, const char *k);

int mlua_buffersize();

void luaopen_mlua(lua_State *L);
