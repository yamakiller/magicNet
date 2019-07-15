package testing

import (
	"fmt"
	"magicNet/script"
)

// TestJSStack : xxx
func TestJSStack() {
	jsstack := script.NewJSStack()
	//jslib.JSRuntimeBundle(jsstack)
	/*jsstack.SetFunc("Print", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})*/

	val, err := jsstack.ExecuteScriptFile("js/test.js")
	fmt.Printf("test.js=%v,%v", val, err)
}
