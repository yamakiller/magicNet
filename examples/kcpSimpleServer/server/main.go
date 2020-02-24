package main

import (
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/yamakiller/magicLibs/encryption/dh64"
	"github.com/yamakiller/magicLibs/net/middle"

	_ "github.com/mkevac/debugcharts"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/log"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/examples/kcpSimpleServer/server/ado"
	"github.com/yamakiller/magicNet/netboxs"
	"github.com/yamakiller/magicNet/netmsgs"
)

var (
	kcpBox *netboxs.KCPBox
)

func onAccept(context *boxs.Context) {
	request := context.Message().(*netmsgs.Accept)
	kcpBox.OpenTo(request.Sock)
	context.Info("accept connect socket %d", request.Sock)
}

func onMessage(context *boxs.Context) {
	request := context.Message().(*netmsgs.Message)
	context.Info("socket %d message: %+v", request.Sock, request.Data)
}

func onClosed(context *boxs.Context) {
	//request := context.Message().(*netmsgs.Closed)
	//context.Info("closed connect socket %d", request.Sock)
}

func onError(context *boxs.Context) {
	request := context.Message().(*netmsgs.Error)
	context.Info("closed connect socket %d error %+s", request.Sock, request.Err.Error())
}

func main() {
	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	hlog := logrus.New()
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "white+h",
		TimestampStyle: "black+h"})
	hlog.SetFormatter(formatter)
	hlog.SetOutput(os.Stdout)

	logSystem := &log.DefaultAgent{}
	logSystem.WithHandle(hlog)

	engine := actors.New(nil)
	engine.WithLogger(logSystem)

	kcpSrv := &netboxs.KCPBox{
		Box:           *boxs.SpawnBox(nil),
		RecvWndSize:   128,
		SendWndSize:   128,
		RecvQueueSize: 32,
		NoDelay:       1,
		Interval:      10,
		Resend:        2,
		Nc:            1,
		RxMinRto:      10,
		FastResend:    1,
		Middleware: &ado.TestMiddleServe{
			SnkMiddleServe: *middle.SpawnSnkMiddleServe(12001,
				dh64.DefaultP,
				dh64.DefaultG,
				time.Second),
		},
	}
	kcpSrv.WithPool(&ado.ConnPools{})

	_, err := engine.New(func(pid *actors.PID) actors.Actor {
		kcpSrv.WithPID(pid)
		kcpSrv.WithMax(1024)
		kcpSrv.Register(reflect.TypeOf(&netmsgs.Accept{}), onAccept)
		kcpSrv.Register(reflect.TypeOf(&netmsgs.Message{}), onMessage)
		kcpSrv.Register(reflect.TypeOf(&netmsgs.Closed{}), onClosed)
		kcpSrv.Register(reflect.TypeOf(&netmsgs.Error{}), onError)
		return kcpSrv
	})

	watch := util.SignalWatch{}
	closed := make(chan bool)

	if err != nil {
		logSystem.Error("", "启动TCPBoxs失败")
		goto exit
	}
	if err = kcpSrv.ListenAndServe("0.0.0.0:12000"); err != nil {
		logSystem.Error("", "监听失败, %s", err.Error())
		goto exit
	}

	logSystem.Info("", "监听成功")

	watch.Initial(func() {
		closed <- true
	})

	watch.Watch()

	for {
		select {
		case <-closed:
		}
		break
	}

exit:
	logSystem.Info("", "退出系统...")
	kcpSrv.ShutdownWait()
	logSystem.Close()

}
