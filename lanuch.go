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

	//LuaReg测试-OK--------------------------------
	/*af := make([]mlua.LuaReg, 1, 1)
	af[0].Name = "test2"
	af[0].Func = test2

	fmt.Printf("start :%d\n", L.GetTop())

	L.NewTable()
	//L.PushInteger(1)
	fmt.Printf("start 1:%d\n", L.GetTop())
	L.SetFuncs(af, 0)

	fmt.Printf("start 2:%d\n", L.GetTop())
	L.SetGlobal("mm")

	L.DoString("print('ooooo') mm.test2()")*/

	var ispass bool
	fmt.Scanln(&ispass)
}
