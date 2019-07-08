package main

import (
	"magicNet/engine/evtchan"
	"magicNet/test"
)

func main() {
	test.TestEmpty()
	evtchan.TestGlobalEventChan()
	//test.TestLanuchHook()
	//test.TestLanuchHookDefault()
	test.TestWait()
}
