package testing

import (
	"fmt"
	"os"
)

// TestDir ：测试目录
func TestDir() {
	dir, err := os.Getwd()

	fmt.Println("uuuuu:", dir, err)

}
