package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//SignalWatch doc
//@Struct SignalWatch @Summary signal watch proccesser
type SignalWatch struct {
	_c chan os.Signal
	_e sync.WaitGroup
	_f func()
}

//Initial doc
//@Method Initial @Summary Initialization signal watcher
//@Param (func()) Signal response back call function
func (slf *SignalWatch) Initial(f func()) {
	slf._f = f
	slf._c = make(chan os.Signal)
	slf._e = sync.WaitGroup{}
	signal.Notify(slf._c, os.Interrupt, os.Kill, syscall.SIGTERM)
}

//Watch doc
//@Method Watch @Summary start watch signal
func (slf *SignalWatch) Watch() {
	slf._e.Add(1)
	go func() {
		defer slf._e.Done()
		for s := range slf._c {
			switch s {
			case os.Interrupt, os.Kill, syscall.SIGTERM:
				slf._f()
				return
			default:
				break
			}
		}
	}()
}

//Wait doc
//@Method Wait @Summary wait signal watcher exit
func (slf *SignalWatch) Wait() {
	slf._e.Wait()
}
