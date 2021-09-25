package util

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/takama/daemon"
)

//var dependencies = []string{"dummy.service"}

//Application Daemon 服务对象接口
type Application interface {
	Name() string
	Desc() string
	Open() (string, error)
	Close()
}

//SpawnDaemon 创建一个Daemon
func SpawnDaemon(app Application) *Daemon {
	d, err := daemon.New(app.Name(),
		app.Desc(),
		daemon.SystemDaemon)

	if err != nil {
		return nil
	}

	return &Daemon{
		_d:   d,
		_s:   make(chan os.Signal, 1),
		_app: app,
	}
}

//Daemon doc
type Daemon struct {
	_d   daemon.Daemon
	_s   chan os.Signal
	_app Application
}

//Open 打开Daemon
func (slf *Daemon) Open() (string, error) {
	//TODO: 需要增加，对core文件的收集
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			if len(os.Args[2:]) == 0 {
				return slf._d.Install()
			}
			return slf._d.Install(os.Args[2:]...)
		case "remove":
			return slf._d.Remove()
		case "start":
			return slf._d.Start()
		case "stop":
			return slf._d.Stop()
		case "status":
			return slf._d.Status()
		case "help":
			usage := "Usage: " + slf._app.Name() + " install | remove | start | stop | status"
			return usage, nil
		default:
		}
	}

	signal.Notify(slf._s, os.Interrupt, os.Kill, syscall.SIGTERM)
	s, err := slf._app.Open()
	if err != nil {
		return s, err
	}

	for {
		select {
		case killSignal := <-slf._s:
			slf._app.Close()
			if killSignal == os.Interrupt {
				return "was interruped by system signal", nil
			}

			return "was killed", nil
		}
	}
}
