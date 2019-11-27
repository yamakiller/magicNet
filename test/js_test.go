package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicLibs/script"
)

// TestJSStack : xxx
func TestJSStack(t *testing.T) {
	jsstack := script.NewJSStack()
	//jslib.JSRuntimeBundle(jsstack)
	/*jsstack.SetFunc("Print", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})*/

	val, err := jsstack.ExecuteScriptFile("js/test.js")
	fmt.Printf("test.js=%v,%v", val, err)
}
