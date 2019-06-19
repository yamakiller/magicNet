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
	L.OpenLibs()

	L.Register("test2", test2)

	var ispass bool
	fmt.Scanln(&ispass)
}
