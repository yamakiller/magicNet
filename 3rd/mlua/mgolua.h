#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"
#include <stdint.h>

typedef int (*lua_GOWrapperCaller) (lua_State *L, unsigned int wrapperid, int top);

int mlua_get_lib_version();

lua_State* mlua_newstate(void* goallocf);

void mlua_setgostate(lua_State *L, size_t gohandle);

size_t mlua_getgostate(lua_State* L);

int mlua_loadfile(lua_State *L, const char *filename);

void mlua_push_go_wrapper(lua_State* L, unsigned int wrapperid);

int mlua_pcall(lua_State* L, int nargs, int nresults, int errfunc);

int mlua_tointeger (lua_State *L, int idx);

void luaopen_mlua(lua_State *L);
