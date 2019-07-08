package evtchan

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
		fmt.Printf("测试一下:%v\n", v)
		return
	}
	fmt.Printf("虽然执行了但并未成功解析数据\n")
}

// TestGlobalEventChan  ：全局 通道测试函数
func TestGlobalEventChan() {
	Subscribe(testGolbalSubRun)
	Publish(&TestLetterEvent{0, 1, 2})
}
