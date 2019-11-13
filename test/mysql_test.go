package test

import (
	"fmt"
	"testing"

	"github.com/yamakiller/magicNet/library"
)

func TestMySql(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/testGo"
	handle := library.MySQLDB{}
	if err := handle.Init(dsn, 10, 1, 3600*1000); err != nil {
		fmt.Println("connect mysql fail:", err)
		return
	}

	result, err := handle.Query("select * from test1")
	if err != nil {
		fmt.Println("sql opation fail:", err)
		handle.Close()
		return
	}

	if result == nil || result.GetRow() == 0 {
		handle.Close()
		return
	}

	fmt.Println("no data:", result.GetRow())
	for result.HashNext() {
		result.Next()
		a, _ := result.GetAsNameValue("c1")
		fmt.Println(a.ToInt32())
		b, _ := result.GetAsNameValue("c2")
		fmt.Println(b.ToInt32())
		c, _ := result.GetAsNameValue("c3")
		fmt.Println(c.ToFloat())
		d, _ := result.GetAsNameValue("c4")
		fmt.Println(d.ToDateTime())
		e, _ := result.GetAsNameValue("c5")
		fmt.Println(e.ToTimeStamp())
		f, _ := result.GetAsNameValue("c6")
		fmt.Println(f.ToString())
	}

	result.Close()
	handle.Close()
}
