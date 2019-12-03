package test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yamakiller/magicLibs/logger"
	"github.com/yamakiller/magicNet/engine/actor"
)

var wait sync.WaitGroup

type myMessage struct {
	i   int32
	pid *actor.PID
}

type routerActor struct{}
type tellerActor struct{}

func (state *routerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		//logger.Info(context.Self().ID, "处理一个myMessage %08x， %d", msg.pid.ID, msg.i)
		atomic.AddInt32(&msg.i, 1)
		wait.Done()
	case *actor.Started:
		logger.Error(context.Self().ID, "router Actor 已启动\n")
	default:
		logger.Error(context.Self().ID, "其他消息\n")
	}
}

func (state *tellerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		for i := 0; i < 100; i++ {
			context.Send(msg.pid, msg)
			time.Sleep(10 * time.Millisecond)
		}

	}
}

//TestActorContext desc
//@Method TestActorContext desc: test actor context
func TestActorContext(t *testing.T) {
	schedulerContext := actor.DefaultSchedulerContext
	wait.Add(100 * 1000)

	/*schedulerContext.SetSenderMiddleware(func(_ actor.SenderContext, target *actor.PID, pack *actor.MessagePack) {
				target.set
	})*/

	tmp := &actor.Agnets{}
	rpid := schedulerContext.Make(tmp.SetMakeActor(func() actor.Actor { return &routerActor{} }))
	agnet := actor.AgnetFromActorMaker(func() actor.Actor { return &tellerActor{} })

	for i := 0; i < 1000; i++ {
		pid := schedulerContext.Make(agnet)
		schedulerContext.Send(pid, &myMessage{int32(i), rpid})
	}

	wait.Wait()
}
