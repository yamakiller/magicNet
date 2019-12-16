package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicNet/handler/implement/buffer"
)

func TestRingBuffer(t *testing.T) {

	bfw := "test001-test002-test003"
	bfw2 := "ubkf109fcts"
	bfl := len([]byte(bfw))
	bf1 := buffer.NewRing(128)
	for {
		n, _ := bf1.Write([]byte(bfw))
		if n < bfl {
			fmt.Printf("write complate need:%d,%d\n", bfl, n)
			break
		}
	}

	bf1.Truncated(10)

	i := 0
	var need int
	for {
		need = bfl
		if bf1.Len() == 0 {
			break
		}

		if need > bf1.Len() {
			need = bf1.Len()
		}

		d := bf1.Read(need)
		fmt.Printf("read %s\n", string(d))

		if i == 0 {
			bf1.Write([]byte(bfw2))
		}
		i++
	}
}
