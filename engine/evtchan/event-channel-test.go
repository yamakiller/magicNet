package evtchan

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月08日 13:52:48
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月08日 15:29:14
 * @Description: 事件通道测试代码
 */

import (
	"fmt"
)

// TestLetterEvent ： 测试事件
type TestLetterEvent struct {
	PID     uint32      // The invalid process, to which the message was sent
	Message interface{} // The message that could not be delivered
	Sender  uint32      // the process that sent the Message
}

// TestGolbalSubRun ： 全局订阅测试函数
func testGolbalSubRun(evt interface{}) {
	if v, ok := evt.(*TestLetterEvent); ok {
		fmt.Printf("测试一下:PID:%d, Sender:%d\n", v.PID, v.Sender)
		return
	}
	fmt.Printf("虽然执行了但并未成功解析数据\n")
}

// TestGlobalEventChan  ：全局 通道测试函数
func TestGlobalEventChan() {
	Subscribe(testGolbalSubRun)
	for i := 0; i < 100; i++ {
		testi := i + 1
		Publish(&TestLetterEvent{uint32(i), 0, uint32(testi)})
	}
}
