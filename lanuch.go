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

type act struct {
	A int
	B int
	C int
}


func main() {
	L := mlua.NewState()
	L.OpenLibs()

  //结构及 函数注册测试完成---------------------------
	/*td := &act{1, 2, 3}

	L.PushUserGoStruct(td)

	var bbb act
  L.ToUserGoStruct(-1, &bbb)

	//fmt.Print("as:", unsafe.Sizeof(td))

	//tm := (*act)(unsafe.Pointer(L.ToGoStruct(-1)))
	fmt.Print( td.A, bbb.A)*/

	//L.Register("test2", test2)
  //----------------------------------------------


	var ispass bool
	fmt.Scanln(&ispass)
}
