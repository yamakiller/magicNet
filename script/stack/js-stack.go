package stack

import (
	"errors"

	"github.com/yamakiller/magicNet/engine/files"

	"github.com/robertkrimen/otto"
)

var (
	// ErrJSNotFindFile :
	ErrJSNotFindFile = errors.New("script file does not exist")
	// ErrJSNotFileData :
	ErrJSNotFileData = errors.New("did not get file data")
)

// JSStack : javascirpt 虚拟器
type JSStack struct {
	state *otto.Otto
}

// MakeJSStack : 制作 JS 虚拟机
func MakeJSStack() *JSStack {
	return &JSStack{otto.New()}
}

// SetInt : 给JS脚本设置Int变量
func (slf *JSStack) SetInt(name string, val int) {
	slf.state.Set(name, val)
}

// SetFloat : 给JS脚本设置Float 32 变量
func (slf *JSStack) SetFloat(name string, val float32) {
	slf.state.Set(name, val)
}

// SetDouble :  给JS脚本设置Float 64 变量
func (slf *JSStack) SetDouble(name string, val float64) {
	slf.state.Set(name, val)
}

// SetBoolean : 给JS脚本设置Bool 变量
func (slf *JSStack) SetBoolean(name string, val bool) {
	slf.state.Set(name, val)
}

// SetString : 给JS脚本设置String变量
func (slf *JSStack) SetString(name string, val string) {
	slf.state.Set(name, val)
}

// SetFunc : 设置js脚本调用Go的函数
func (slf *JSStack) SetFunc(name string, fun interface{}) {
	slf.state.Set(name, fun)
}

// ExecuteScriptFile : 执行脚本文件
func (slf *JSStack) ExecuteScriptFile(filename string) (otto.Value, error) {
	tmpFileName := files.GetFullPathForFilename(filename)
	if !files.IsFileExist(tmpFileName) {
		return otto.Value{}, ErrJSNotFindFile
	}

	data := files.GetDataFromFile(tmpFileName)
	if data.IsNil() {
		return otto.Value{}, ErrJSNotFileData
	}

	return slf.state.Run(string(data.GetBytes()))
}
