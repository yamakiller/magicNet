package main

import (
	"net/http"
	"os"
	"reflect"

	_ "github.com/mkevac/debugcharts"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/log"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/examples/tcpSimpleServer/server/ado"
	"github.com/yamakiller/magicNet/netboxs"
	"github.com/yamakiller/magicNet/netmsgs"
)

var (
	tcpBox *netboxs.TCPBox
)

func onAccept(context *boxs.Context) {
	request := context.Message().(*netmsgs.Accept)
	tcpBox.OpenTo(request.Sock)
	context.Info("accept connect socket %d", request.Sock)
}

func onMessage(context *boxs.Context) {
	request := context.Message().(*netmsgs.Message)
	context.Info("socket %d message: %+v", request.Sock, request.Data)
}

func onClosed(context *boxs.Context) {
	request := context.Message().(*netmsgs.Closed)
	context.Info("closed connect socket %d", request.Sock)
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

	tcpService, _ := netboxs.Spawn(netboxs.ModeTCPListener, &ado.ConnPools{})
	_, err := engine.New(func(pid *actors.PID) actors.Actor {
		tcpService.(*netboxs.TCPBox).WithPID(pid)
		tcpService.(*netboxs.TCPBox).WithMax(1024)
		tcpService.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Accept{}), onAccept)
		tcpService.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Message{}), onMessage)
		tcpService.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Closed{}), onClosed)
		return tcpService
	})

	watch := util.SignalWatch{}
	closed := make(chan bool)

	if err != nil {
		logSystem.Error("", "启动TCPBoxs失败")
		goto exit
	}

	tcpBox = tcpService.(*netboxs.TCPBox)
	if err = tcpService.(*netboxs.TCPBox).ListenAndServe("0.0.0.0:12000"); err != nil {
		logSystem.Error("", "监听失败, %s", err.Error())
		tcpBox = nil
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
	tcpService.(*netboxs.TCPBox).ShutdownWait()
	logSystem.Close()
}
