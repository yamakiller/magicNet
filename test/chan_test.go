package test

import (
	"fmt"
	"testing"
)

//TestChan @Summary
//@Method TestChan @Summary test channge
func TestChan(t *testing.T) {
	u := make(chan int)
	go func() {
		a := <-u
		fmt.Println(a)
	}()

	close(u)
}
