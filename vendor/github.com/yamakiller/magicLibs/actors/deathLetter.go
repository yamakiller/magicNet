package actors

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/yamakiller/magicLibs/actors/messages"
)

//Death actor 消息
type Death struct {
	PID     *PID
	Message interface{}
	Sender  *PID
}

type deathLetter struct {
	_parent *Core
	_sub    chan interface{}
	_closed chan bool
	_wait   sync.WaitGroup
}

func (slf *deathLetter) run() {

	defer func() {
		close(slf._closed)
		slf._wait.Done()
	}()

	var msg interface{}
	for {
		select {
		case <-slf._closed:
			goto exit
		case msg = <-slf._sub:
			if deathLetter, ok := msg.(*Death); ok {
				if deathLetter.Sender != nil {
					slf._parent._log.Error(fmt.Sprintf("[%s]", deathLetter.Sender.ToString()), "Death Dest PID :%s Message:%+v", deathLetter.PID.ToString(), deathLetter.Message)
				} else {
					slf._parent._log.Error("", "Death Dest PID: %s Message:%+v", deathLetter.PID.ToString(), reflect.TypeOf(deathLetter.Message))
				}
			}
		}
	}
exit:
}

func (slf *deathLetter) close() {
	for {
		if len(slf._sub) <= 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	slf._closed <- true
	slf._wait.Wait()
}

func (slf *deathLetter) overloadUsrMessage() int {
	return 0
}

func (slf *deathLetter) postUsrMessage(pid *PID, message interface{}) {
	_, msg, sender := UnWrapPack(message)
	slf._sub <- &Death{
		PID:     pid,
		Message: msg,
		Sender:  sender,
	}
}

func (slf *deathLetter) postSysMessage(pid *PID, message interface{}) {
	slf._sub <- &Death{
		PID:     pid,
		Message: message,
	}
}

func (slf *deathLetter) Stop(pid *PID) {
	slf.postSysMessage(pid, &messages.Stop{})
}
