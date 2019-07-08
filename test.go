package main

import (
	"magicNet/engine/evtchan"
	"magicNet/testing"
)

func main() {
	testing.TestEmpty()
	evtchan.TestGlobalEventChan()
	//test.TestLanuchHook()
	//test.TestLanuchHookDefault()
	//testing.TestWait()
}
