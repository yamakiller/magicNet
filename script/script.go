package script

import (
	"github.com/yamakiller/magicNet/script/stack"

	"github.com/yamakiller/magicNet/script/jslib"
)

// NewJSStack : 创建爱你一个js虚拟机
func NewJSStack() *stack.JSStack {
	stack := stack.MakeJSStack()
	jslib.Bundle(stack)
	return stack
}
