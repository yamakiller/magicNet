package testing

import "fmt"

//TestChan : 管道
func TestChan() {
	u := make(chan int)
	go func() {
		a := <-u
		fmt.Println(a)
	}()

	close(u)
}
