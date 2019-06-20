#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include <stdint.h>

#define GOLUA_PANIC_MSG_WARAPPER "golua_panicmsg_warapper"


typedef int (*lua_GOWrapperCaller) (lua_State *L, unsigned int wrapperid, int top);

int mlua_get_lib_version();

lua_State* mlua_newstate(void* goallocf);

void mlua_setallocf(lua_State* L, void* goallocf);

void mlua_setgostate(lua_State *L, void *goluaState);

void* mlua_getgostate(lua_State* L);

int mlua_loadfile(lua_State *L, const char *filename);

int mlua_loadbuffer(lua_State *L, const char *buffer, size_t sz, const char* name);

void mlua_push_go_wrapper(lua_State* L, unsigned int wrapperid);

int mlua_pcall(lua_State* L, int nargs, int nresults, int errfunc);

lua_Integer mlua_tointeger(lua_State *L, int idx);

lua_Number mlua_tonumber(lua_State *L, int idx);

const char *mlua_tostring(lua_State *L, int idx);

int mlua_error(lua_State *L, const char *fmt);

void mlua_replace(lua_State *L, int idx);

void mlua_pushglobaltable(lua_State *L);

void mlua_pushgostruct(lua_State *L, unsigned int wrapperid);

unsigned int mlua_isgostruct(lua_State *, int idx);

int mlua_getmetatable(lua_State *L, const char *k);

int mlua_buffersize();

void luaopen_mlua(lua_State *L);
