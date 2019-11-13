package test

import (
	"fmt"
	"os"
	"testing"
)

//TestDir Desc:
//@method TestDir desc: test Getwd()
func TestDir(t *testing.T) {
	dir, err := os.Getwd()

	fmt.Println("uuuuu:", dir, err)

}

func TestRectPoint(t *testing.T) {

}
