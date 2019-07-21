package jslib

import (
	"github.com/yamakiller/magicNet/script/stack"

	"github.com/robertkrimen/otto"
)

// Bundle : 基础绑定 js库，在此处扩展
func Bundle(stack *stack.JSStack) {
	stack.SetFunc("Refer", refer)
}

func refer(js otto.FunctionCall) otto.Value {
	switch js.Argument(0).String() {
	case "runtime":
		return jsruntimeBundle(js)
	default:
		return otto.Value{}
	}
}
