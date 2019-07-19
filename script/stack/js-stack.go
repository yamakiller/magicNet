package stack

import (
	"errors"
	"magicNet/engine/files"

	"github.com/robertkrimen/otto"
)

var (
	ErrJSNotFindFile = errors.New("script file does not exist")
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
func (js *JSStack) SetInt(name string, val int) {
	js.state.Set(name, val)
}

// SetFloat : 给JS脚本设置Float 32 变量
func (js *JSStack) SetFloat(name string, val float32) {
	js.state.Set(name, val)
}

// SetDouble :  给JS脚本设置Float 64 变量
func (js *JSStack) SetDouble(name string, val float64) {
	js.state.Set(name, val)
}

// SetBoolean : 给JS脚本设置Bool 变量
func (js *JSStack) SetBoolean(name string, val bool) {
	js.state.Set(name, val)
}

// SetString : 给JS脚本设置String变量
func (js *JSStack) SetString(name string, val string) {
	js.state.Set(name, val)
}

// SetFunc : 设置js脚本调用Go的函数
func (js *JSStack) SetFunc(name string, fun interface{}) {
	js.state.Set(name, fun)
}

// ExecuteScriptFile : 执行脚本文件
func (js *JSStack) ExecuteScriptFile(filename string) (otto.Value, error) {
	tmpFileName := files.GetFullPathForFilename(filename)
	if !files.IsFileExist(tmpFileName) {
		return otto.Value{}, ErrJSNotFindFile
	}

	data := files.GetDataFromFile(tmpFileName)
	if data.IsNil() {
		return otto.Value{}, ErrJSNotFileData
	}

	return js.state.Run(string(data.GetBytes()))
}
