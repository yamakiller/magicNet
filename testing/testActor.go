package testing

import (
	fmt "fmt"
	"magicNet/engine/actor"
	"magicNet/engine/logger"
	"sync"
	"sync/atomic"
	"time"
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
		fmt.Printf("%v 处理一个myMessage %v \n", context.Self(), msg)
		atomic.AddInt32(&msg.i, 1)
		wait.Done()
	case *actor.Started:
		logger.Error(context.Self().ID, "router Actor 已启动\n")
	}
}

func (state *tellerActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *myMessage:
		for i := 0; i < 1; i++ {
			logger.Info(context.Self().ID, "发送数据")
			context.Send(msg.pid, msg)
			time.Sleep(10 * time.Millisecond)
		}

	}
}

// TestActorContext : 测试大量Actor 消息发送
func TestActorContext() {
	schedulerContext := actor.DefaultSchedulerContext
	wait.Add(1 * 1)

	/*schedulerContext.SetSenderMiddleware(func(_ actor.SenderContext, target *actor.PID, pack *actor.MessagePack) {
				target.set
	})*/

	tmp := &actor.Agnets{}
	rpid := schedulerContext.Make(tmp.SetMakeActor(func() actor.Actor { return &routerActor{} }))

	if rpid == nil {

	}
	//schedulerContext.Send(rpid, &myMessage{int32(1), rpid})

	/*agnet := actor.AgnetFromActorMaker(func() actor.Actor { return &tellerActor{} })

	for i := 0; i < 1; i++ {
		pid := schedulerContext.Make(agnet)
		schedulerContext.Send(pid, &myMessage{int32(i), rpid})
	}*/

	//wait.Wait()
}
