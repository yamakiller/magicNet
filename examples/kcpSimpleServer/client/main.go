package main

import (
	"fmt"
	"time"

	"github.com/yamakiller/magicLibs/args"
	"github.com/yamakiller/magicNet/examples/kcpSimpleServer/client/ado"
)

func main() {
	args.Instance().Parse()
	addr := args.Instance().GetString("-p", "127.0.0.1:12000")
	maxConn := args.Instance().GetInt("-n", 1)
	deply := args.Instance().GetInt("-d", 300)
	checkNum := args.Instance().GetInt("-c", 1)
	timeOut := args.Instance().GetInt("-t", 1000)

	var connected int
	var sendCount int
	var failCount int
	var clients []*ado.Clt
	timeDeply := time.Duration(deply)

	for i := 0; i < maxConn; i++ {
		c := ado.NewClt(uint32(i + 1))

		fmt.Print("开始连接:", addr)
		if err := c.Connect(addr, time.Duration(timeOut)*time.Millisecond); err != nil {
			c.Close()
			fmt.Println(" 连接失败 ", err)
			continue
		}

		fmt.Println(" 连接成功")
		go c.ReadLoop()

		c.Timeout = time.Now()
		clients = append(clients, c)
		connected++
	}

	for {
		curtime := time.Now()
		for i := 0; i < len(clients); {
			cc := clients[i]
			diff := curtime.Sub(cc.Timeout)
			if diff.Milliseconds() >= int64(timeDeply) {
				if cc.Check >= checkNum {
					clients = append(clients[0:i], clients[i+1:]...)
					cc.Close()
					continue
				}

				if err := cc.SendTo("abcd"); err != nil {
					failCount++
					clients = append(clients[0:i], clients[i+1:]...)
					cc.Close()
					continue
				}

				sendCount++
				cc.Timeout = curtime
				cc.Check++
			}
			i++
		}

		if len(clients) == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("总连接次数:", maxConn)
	fmt.Println("完成连接数:", connected)
	fmt.Println("发送成功次数:", sendCount)
	fmt.Println("发送失败次数:", failCount)
}
