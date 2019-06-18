package main

import (
	"fmt"
	"magicNet/3rd/mlua"
)

func test2(L *mlua.State) int {
	fmt.Print("ssssssss\n")
	return 0
}

func test3(L *mlua.State) int {
	fmt.Print("cccccc\n")
	return 0
}

func main() {
	L := mlua.NewState()
	defer L.Close()
	L.OpenLibs()
	fmt.Print("ddddddddd\n")
	L.Register("test2", test2)
	L.Register("test3", test3)

	//L.PushGoFunction(test2)
	L.DoString("test2() test3() test2()")
	fmt.Print("end\n")
}
