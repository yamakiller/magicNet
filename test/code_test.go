package test

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

//TestDir Desc:
//@Method TestDir desc: test Getwd()
func TestDir(t *testing.T) {
	dir, err := os.Getwd()

	fmt.Println("uuuuu:", dir, err)
	fmt.Println("system:", runtime.GOOS)
}

func TestRectPoint(t *testing.T) {

}
