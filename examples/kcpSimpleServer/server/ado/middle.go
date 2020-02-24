package ado

import (
	"fmt"

	"github.com/yamakiller/magicLibs/net/middle"
)

//TestMiddleServe 测试中间件
type TestMiddleServe struct {
	middle.SnkMiddleServe
}

//Error ...
func (slf *TestMiddleServe) Error(err error) {
	fmt.Println("Error:", err)
}

//Debug ...
func (slf *TestMiddleServe) Debug(err error) {
	fmt.Println("Debug:", err)
}
