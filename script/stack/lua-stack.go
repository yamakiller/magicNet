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
func (slf *LuaStack) GetLuaState() *mlua.State {
	return slf._l
}

//OpenLibs  library
func (slf *LuaStack) OpenLibs() {
	slf._l.OpenLibs()
}

// AddSreachPath : 添加LUA搜索路径
func (slf *LuaStack) AddSreachPath(path string) {
	slf._l.GetGlobal("package")
	slf._l.GetField(-1, "path")
	curPath := slf._l.ToString(-1)
	newPath := fmt.Sprintf("%s;%s/?.lua", curPath, path)
	slf._l.PushString(newPath)
	slf._l.SetField(-3, "path")
	slf._l.Pop(2)
}

// AddLuaLoader : LUA载入器
func (slf *LuaStack) AddLuaLoader(f *mlua.LuaGoFunction) {
	if f == nil {
		return
	}
	slf._l.GetGlobal("package")
	slf._l.GetField(-1, "preload")

	slf._l.PushGoFunction(*f)
	for i := slf._l.RawLen(-2) + 1; i > 2; {
		slf._l.RawGetI(-2, int64(i-1))
		slf._l.RawGetI(-3, int64(i))
		i--
	}
	slf._l.RawSetI(-2, 2)

	slf._l.SetField(-2, "preload")
	slf._l.Pop(1)
}

// ExecuteFunction : 执行LUA 函数
func (slf *LuaStack) ExecuteFunction(numArgs int) (int, error) {
	funcIndex := -(numArgs + 1)
	if !slf._l.IsFunction(funcIndex) {
		slf._l.Pop(numArgs + 1)
		return 0, nil
	}

	traceback := 0
	slf._l.GetGlobal("__G__TRACKBACK__")
	if !slf._l.IsFunction(-1) {
		slf._l.Pop(1)
	} else {
		slf._l.Insert(funcIndex - 1)
		traceback = funcIndex - 1
	}

	error := slf._l.PCall(numArgs, 1, traceback)
	if error != 0 {
		if traceback == 0 {
			err := slf._l.ToString(-1)
			slf._l.Pop(1)
			return 0, errors.New(err)
		}
		slf._l.Pop(2)
		return 0, errors.New("lua unknown error")
	}

	ret := 0
	if slf._l.IsNumber(-1) {
		ret = int(slf._l.ToInteger(-1))
	} else if slf._l.IsBoolean(-1) {
		if slf._l.ToBoolean(-1) {
			ret = 1
		} else {
			ret = 0
		}
	}

	slf._l.Pop(1)

	if traceback != 0 {
		slf._l.Pop(1)
	}

	return ret, nil
}

// ExecuteString : 执行LUA字符串
func (slf *LuaStack) ExecuteString(codes string) (int, error) {
	if slf._l.LoadString(codes) != 0 {
		err := errors.New(slf._l.ToString(-1))
		slf._l.Pop(1)
		return 0, err
	}
	return slf.ExecuteFunction(0)
}

// ExecuteScriptFile ： 执行LUA脚本文件
func (slf *LuaStack) ExecuteScriptFile(fileName string) (int, error) {
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

	if _, err := slf.loadBuffer(&data.GetBytes()[0], uint(data.GetSize()), tmp); err != nil {
		return 0, err
	}

	return slf.ExecuteFunction(0)
}

// Clean : 清空堆栈
func (slf *LuaStack) Clean() {
	slf._l.SetTop(0)
}

// PushInt : 插入Int
func (slf *LuaStack) PushInt(intValue int) {
	slf._l.PushInteger(int64(intValue))
}

// PushLong : 插入64位Int
func (slf *LuaStack) PushLong(longValue int64) {
	slf._l.PushInteger(longValue)
}

// PushFloat : 插入 float
func (slf *LuaStack) PushFloat(floatValue float32) {
	slf._l.PushNumber(float64(floatValue))
}

// PushDouble : 插入 float64
func (slf *LuaStack) PushDouble(doubleValue float64) {
	slf._l.PushNumber(float64(doubleValue))
}

// PushBoolean : 插入 bool
func (slf *LuaStack) PushBoolean(boolValue bool) {
	slf._l.PushBoolean(boolValue)
}

// PushString : 插入字符串
func (slf *LuaStack) PushString(stringValue string) {
	slf._l.PushString(stringValue)
}

// PushNil : 插入一个 Nil
func (slf *LuaStack) PushNil() {
	slf._l.PushNil()
}

// Register : 注册闭包函数
func (slf *LuaStack) Register(f mlua.LuaGoFunction, name string, args ...interface{}) {
	slf._l.PushGoClosure(f, args...)
	slf._l.SetGlobal(name)
}

// ReLoad : 重新载入
func (slf *LuaStack) ReLoad(moduleFileName string) (int, error) {
	if len(moduleFileName) == 0 {
		return 0, fmt.Errorf("reload %s fail", moduleFileName)
	}

	slf._l.GetGlobal("package")
	slf._l.GetField(-1, "loaded")
	slf._l.PushString(moduleFileName)
	slf._l.GetTable(-2)
	if !slf._l.IsNil(-1) {
		slf._l.PushString(moduleFileName)
		slf._l.PushNil()
		slf._l.SetTable(-4)
	}
	slf._l.Pop(3)

	name := moduleFileName
	require := fmt.Sprintf("require '%s'", name)
	return slf.ExecuteString(require)
}

func (slf *LuaStack) loadBuffer(chunk *byte, chunkSize uint, chunkName string) (int, error) {
	r := slf._l.LoadBuffer(chunk, chunkSize, chunkName)

	if r != 0 {

		err := slf._l.ToString(-1)
		slf._l.Pop(1)

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
func (slf *LuaStack) Shutdown() {
	slf._l.Close()
}
