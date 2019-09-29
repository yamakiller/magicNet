package stack

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yamakiller/magicNet/engine/files"
	"github.com/yamakiller/magicNet/util"

	"github.com/yamakiller/mgolua/mlua"
)

const (
	byteLuaFileExt    = ".luac"
	notByteLuaFileExt = ".lua"
)

// LuaStack : LUA虚拟机 堆
type LuaStack struct {
	_l *mlua.State
}

// NewLuaStack create a lua stack
func NewLuaStack() *LuaStack {
	return &LuaStack{_l: mlua.NewState()}
}

// GetLuaState Get LUA virtual machine C object
func (S *LuaStack) GetLuaState() *mlua.State {
	return S._l
}

//OpenLibs  library
func (S *LuaStack) OpenLibs() {
	S._l.OpenLibs()
}

// AddSreachPath : 添加LUA搜索路径
func (S *LuaStack) AddSreachPath(path string) {
	S._l.GetGlobal("package")
	S._l.GetField(-1, "path")
	curPath := S._l.ToString(-1)
	newPath := fmt.Sprintf("%s;%s/?.lua", curPath, path)
	S._l.PushString(newPath)
	S._l.SetField(-3, "path")
	S._l.Pop(2)
}

// AddLuaLoader : LUA载入器
func (S *LuaStack) AddLuaLoader(f *mlua.LuaGoFunction) {
	if f == nil {
		return
	}
	S._l.GetGlobal("package")
	S._l.GetField(-1, "preload")

	S._l.PushGoFunction(*f)
	for i := S._l.RawLen(-2) + 1; i > 2; {
		S._l.RawGetI(-2, int64(i-1))
		S._l.RawGetI(-3, int64(i))
		i--
	}
	S._l.RawSetI(-2, 2)

	S._l.SetField(-2, "preload")
	S._l.Pop(1)
}

// ExecuteFunction : 执行LUA 函数
func (S *LuaStack) ExecuteFunction(numArgs int) (int, error) {
	funcIndex := -(numArgs + 1)
	if !S._l.IsFunction(funcIndex) {
		S._l.Pop(numArgs + 1)
		return 0, nil
	}

	traceback := 0
	S._l.GetGlobal("__G__TRACKBACK__")
	if !S._l.IsFunction(-1) {
		S._l.Pop(1)
	} else {
		S._l.Insert(funcIndex - 1)
		traceback = funcIndex - 1
	}

	error := S._l.PCall(numArgs, 1, traceback)
	if error != 0 {
		if traceback == 0 {
			err := S._l.ToString(-1)
			S._l.Pop(1)
			return 0, errors.New(err)
		}
		S._l.Pop(2)
		return 0, errors.New("lua unknown error")
	}

	ret := 0
	if S._l.IsNumber(-1) {
		ret = int(S._l.ToInteger(-1))
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

	return ret, nil
}

// ExecuteString : 执行LUA字符串
func (S *LuaStack) ExecuteString(codes string) (int, error) {
	if S._l.LoadString(codes) != 0 {
		err := errors.New(S._l.ToString(-1))
		S._l.Pop(1)
		return 0, err
	}
	return S.ExecuteFunction(0)
}

// ExecuteScriptFile ： 执行LUA脚本文件
func (S *LuaStack) ExecuteScriptFile(fileName string) (int, error) {
	tmp := fileName
	pos := strings.LastIndex(tmp, byteLuaFileExt)
	if pos != -1 {
		tmp = util.SubStr(tmp, 0, pos)
	} else {
		pos = strings.LastIndex(tmp, notByteLuaFileExt)
		if pos == (len(tmp) - len(notByteLuaFileExt)) {
			tmp = util.SubStr(tmp, 0, pos)
		}
	}

	tmpFileName := tmp + byteLuaFileExt
	tmpFileName = files.GetFullPathForFilename(tmpFileName)
	if files.IsFileExist(tmpFileName) {
		tmp = tmpFileName
	} else {
		tmpFileName = tmp + notByteLuaFileExt
		tmpFileName = files.GetFullPathForFilename(tmpFileName)
		if !files.IsFileExist(tmpFileName) {
			return 0, fmt.Errorf("cannot open %s:No such file or directory", tmpFileName)
		}
		tmp = tmpFileName
	}

	data := files.GetDataFromFile(tmp)
	if data.IsNil() {
		return 0, fmt.Errorf("%s script not executed correctly", tmp)
	}

	if _, err := S.loadBuffer(&data.GetBytes()[0], uint(data.GetSize()), tmp); err != nil {
		return 0, err
	}

	return S.ExecuteFunction(0)
}

// Clean : 清空堆栈
func (S *LuaStack) Clean() {
	S._l.SetTop(0)
}

// PushInt : 插入Int
func (S *LuaStack) PushInt(intValue int) {
	S._l.PushInteger(int64(intValue))
}

// PushLong : 插入64位Int
func (S *LuaStack) PushLong(longValue int64) {
	S._l.PushInteger(longValue)
}

// PushFloat : 插入 float
func (S *LuaStack) PushFloat(floatValue float32) {
	S._l.PushNumber(float64(floatValue))
}

// PushDouble : 插入 float64
func (S *LuaStack) PushDouble(doubleValue float64) {
	S._l.PushNumber(float64(doubleValue))
}

// PushBoolean : 插入 bool
func (S *LuaStack) PushBoolean(boolValue bool) {
	S._l.PushBoolean(boolValue)
}

// PushString : 插入字符串
func (S *LuaStack) PushString(stringValue string) {
	S._l.PushString(stringValue)
}

// PushNil : 插入一个 Nil
func (S *LuaStack) PushNil() {
	S._l.PushNil()
}

// Register : 注册闭包函数
func (S *LuaStack) Register(f mlua.LuaGoFunction, name string, args ...interface{}) {
	S._l.PushGoClosure(f, args...)
	S._l.SetGlobal(name)
}

// ReLoad : 重新载入
func (S *LuaStack) ReLoad(moduleFileName string) (int, error) {
	if len(moduleFileName) == 0 {
		return 0, fmt.Errorf("reload %s fail", moduleFileName)
	}

	S._l.GetGlobal("package")
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

func (S *LuaStack) loadBuffer(chunk *byte, chunkSize uint, chunkName string) (int, error) {
	r := S._l.LoadBuffer(chunk, chunkSize, chunkName)

	if r != 0 {

		err := S._l.ToString(-1)
		S._l.Pop(1)

		switch r {
		case int(mlua.LUAERRSYNTAX):
			return r, fmt.Errorf("Lua syntax error in buffer %s: %s", chunkName, err)
		case int(mlua.LUAERRMEM):
			return r, fmt.Errorf("Could not load Lua buffer %s", chunkName)
		case int(mlua.LUAERRFILE):
		default:
			break
		}

		return r, errors.New(err)
	}
	return r, nil
}

//Shutdown Close lua_state
func (S *LuaStack) Shutdown() {
	S._l.Close()
}
