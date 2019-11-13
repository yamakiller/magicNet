package test

import (
	"fmt"
	"testing"
)

//TestChan desc:
//@method TestChan desc: test channge
func TestChan(t *testing.T) {
	u := make(chan int)
	go func() {
		a := <-u
		fmt.Println(a)
	}()

	close(u)
}
