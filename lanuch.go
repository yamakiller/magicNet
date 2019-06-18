package main

import (
	"fmt"
	"magicNet/3rd/mlua"
)

func test2(L *mlua.State) int {
	fmt.Print("ssssssss\n")
	return 0
}

func main() {
	L := mlua.NewState()
	defer L.Close()
	L.OpenLibs()
	fmt.Print("ddddddddd\n")
	L.Register("test2", test2)
	//L.PushGoFunction(test2)
	L.DoString("test2()")
	fmt.Print("end\n")
}
