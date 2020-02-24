package ado

import (
	"fmt"

	"github.com/yamakiller/magicLibs/net/middle"
)

//TestMiddleCli 测试中间件
type TestMiddleCli struct {
	middle.SnkMiddleCli
}

//Error ...
func (slf *TestMiddleCli) Error(err error) {
	fmt.Println("Error:", err)
}

//Debug ...
func (slf *TestMiddleCli) Debug(err error) {
	fmt.Println("Debug:", err)
}
