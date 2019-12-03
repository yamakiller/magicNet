package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/core/boot"
	"github.com/yamakiller/magicNet/core/frame"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/handler"
)

var (
	inumber = 0
)

var ppFree = sync.Pool{
	New: func() interface{} {
		inumber++
		s := new(testHandle)
		logger.Info(0, "service:%d", inumber)
		return s
	},
}

type testHandle struct {
	handler.Service
}

//Initial doc
func (slf *testHandle) Initial() {
	slf.Service.Initial()
	slf.RegisterMethod(&actor.Stopped{}, slf.Stoped)
	slf.RegisterMethod(&actor.Terminated{}, slf.Terminated)
}

func (slf *testHandle) Stoped(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogInfo("Stoped")
	slf.Service.Stoped(context, sender, message)
}

func (slf *testHandle) Terminated(context actor.Context, sender *actor.PID, message interface{}) {
	slf.LogInfo("Terminated")
	//slf.Shutdown()
}

type testEngine struct {
	core.DefaultBoot
	core.DefaultService
	core.DefaultWait
}

func (slf *testEngine) InitService() error {
	for i := 0; i < 100; i++ {
		s := handler.Spawn(fmt.Sprintf("service/test/#%d", i+1), func() handler.IService {
			s := ppFree.Get().(*testHandle)
			s.Initial()
			return s
		})

		time.Sleep(time.Duration(100) * time.Millisecond)
		s.Shutdown()
		ppFree.Put(s)

	}
	return nil
}

func TestService(t *testing.T) {
	boot.Launch(func() frame.Framework {
		return &testEngine{}
	})
}
