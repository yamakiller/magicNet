package testing

import (
	"fmt"
	"magicNet/scene/aoi"
	"os"
	"unsafe"
)

// TestDir ：测试目录
func TestDir() {
	dir, err := os.Getwd()

	fmt.Println("uuuuu:", dir, err)

}

type TestRP struct {
	aoi.Rect
	aoi.Point
}

func TestRectPoint() {
	t := TestRP{}
	s := unsafe.Pointer(&t)

	a := (*aoi.Point)(s)

	fmt.Println(a.X)
	/*v := reflect.ValueOf(c).FieldByName("abc")
	if v.IsValid() {
		fmt.Println("1")
	} else {
		fmt.Println("2")
	}*/
}
