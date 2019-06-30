package script

import (
	"fmt"
	"magicNet/files"
	"magicNet/logger"
	"magicNet/util"
	"strings"

	"github.com/yamakiller/mgolua/mlua"
)

const BYTELUA_FILE_EXT = ".luac"
const NOT_BYTELUA_FILE_EXT = ".lua"

type LuaStack struct {
	_l *mlua.State
}

func (S *LuaStack) GetLuaState() *mlua.State {
	return S._l
}

func (S *LuaStack) AddSreachPath(path string) {
	S._l.GetGlobal("searchers")
	S._l.GetField(-1, "path")
	curPath := S._l.ToString(-1)
	newPath := fmt.Sprintf("%s;%s/?.lua", curPath, path)
	S._l.PushString(newPath)
	S._l.SetField(-3, "path")
	S._l.Pop(2)
}

func (S *LuaStack) AddLuaLoader(f *mlua.LuaGoFunction) {
	if f == nil {
		return
	}
	S._l.GetGlobal("searchers")
	S._l.GetField(-1, "loaders")

	S._l.PushGoFunction(*f)
	for i := S._l.RawLen(-2) + 1; i > 2; {
		S._l.RawGetI(-2, int64(i-1))
		S._l.RawGetI(-3, int64(i))
		i--
	}
	S._l.RawSetI(-2, 2)

	S._l.SetField(-2, "loaders")
	S._l.Pop(1)
}

func (S *LuaStack) ExecuteFunction(numArgs int) int {
	funcIndex := -(numArgs + 1)
	if !S._l.IsGFunction(funcIndex) {
		S._l.Pop(numArgs + 1)
		return 0
	}

	traceback := 0
	S._l.GetGlobal("__G__TRACKBACK__")
	if !S._l.IsGFunction(-1) {
		S._l.Pop(1)
	} else {
		S._l.Insert(funcIndex - 1)
		traceback = funcIndex - 1
	}

	error := S._l.PCall(numArgs, 1, traceback)
	if error != 0 {
		if traceback == 0 {
			logger.Error(0, S._l.ToString(-1))
			S._l.Pop(1)
		} else {
			S._l.Pop(2)
		}
		return 0
	}

	ret := 0
	if S._l.IsNumber(-1) {
		ret = S._l.ToInteger(-1)
	} else if S._l.IsBoolean(-1) {
		if S._l.ToBoolean(-1) {
			ret = 1
		} else {
			ret = 0
		}
	}

	S._l.Pop(1)

	if traceback != 0 {
		S._l.Pop(1)
	}
	return ret
}

func (S *LuaStack) ExecuteString(codes string) int {
	S._l.LoadString(codes)
	return S.ExecuteFunction(0)
}

func (S *LuaStack) ExecuteScriptFile(fileName string) int {
	tmp := fileName
	pos := strings.LastIndex(tmp, BYTELUA_FILE_EXT)
	if pos != -1 {
		tmp = util.SubStr(tmp, 0, pos)
	} else {
		pos = strings.LastIndex(tmp, NOT_BYTELUA_FILE_EXT)
		if pos == (len(tmp) - len(NOT_BYTELUA_FILE_EXT)) {
			tmp = util.SubStr(tmp, 0, pos)
		}
	}

	utilFile := files.GetInstance()
	tmpfilename := tmp + BYTELUA_FILE_EXT
	if utilFile.IsFileExist(tmpfilename) {
		tmp = tmpfilename
	} else {
		tmpfilename = tmp + NOT_BYTELUA_FILE_EXT
		if utilFile.IsFileExist(tmpfilename) {
			tmp = tmpfilename
		}
	}

	fullFilePath := utilFile.GetFullPathForFilename(tmp)
	data := utilFile.GetDataFromFile(fullFilePath)
	rn := 0
	if data != nil {
		if S.luaLoadBuffer(data.GetData(), uint(data.GetBytes()), fullFilePath) == 0 {
			rn = S.ExecuteFunction(0)
		}
	}

	return rn
}

func (S *LuaStack) Clean() {
	S._l.SetTop(0)
}

func (S *LuaStack) PushInt(intValue int) {
	S._l.PushInteger(int64(intValue))
}

func (S *LuaStack) PushLong(longValue int64) {
	S._l.PushInteger(longValue)
}

func (S *LuaStack) PushFloat(floatValue float32) {
	S._l.PushNumber(float64(floatValue))
}

func (S *LuaStack) PushDouble(doubleValue float64) {
	S._l.PushNumber(float64(doubleValue))
}

func (S *LuaStack) PushBoolean(boolValue bool) {
	S._l.PushBoolean(boolValue)
}

func (S *LuaStack) PushString(stringValue string) {
	S._l.PushString(stringValue)
}

func (S *LuaStack) PushNil() {
	S._l.PushNil()
}

func (S *LuaStack) ReLoad(moduleFileName string) int {
	if len(moduleFileName) == 0 {
		logger.Error(0, "reload %s fail.", moduleFileName)
		return 1
	}

	S._l.GetGlobal("searchers")
	S._l.GetField(-1, "loaded")
	S._l.PushString(moduleFileName)
	S._l.GetTable(-2)
	if !S._l.IsNil(-1) {
		S._l.PushString(moduleFileName)
		S._l.PushNil()
		S._l.SetTable(-4)
	}
	S._l.Pop(3)

	name := moduleFileName
	require := fmt.Sprintf("require '%s'", name)
	return S.ExecuteString(require)
}

func (S *LuaStack) luaLoadBuffer(chunk *byte, chunkSize uint, chunkName string) int {
	r := S._l.LoadBuffer(chunk, chunkSize, chunkName)

	if r != 0 {
		switch r {
		case int(mlua.LUA_ERRSYNTAX):
			//错误日志
			break
		case int(mlua.LUA_ERRMEM):
			break
		case int(mlua.LUA_ERRFILE):
			break
		default:
			break
		}
	}
	return r
}
