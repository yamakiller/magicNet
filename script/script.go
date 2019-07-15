package script

import (
	"magicNet/script/jslib"
	"magicNet/script/stack"
)

// NewJSStack : 创建爱你一个js虚拟机
func NewJSStack() *stack.JSStack {
	stack := stack.MakeJSStack()
	jslib.Bundle(stack)
	return stack
}
