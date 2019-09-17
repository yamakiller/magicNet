package core

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yamakiller/magicNet/core/version"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/util"
)

const (
	verMajor     = "1"
	verMinor     = "0"
	verPatch     = "0"
	verBuild     = "x64"
	verRevision  = "Beta"
	verCode      = "UTF8"
	verBuildTime = "2019-07-10 16:58"
)

// DefaultCMDLineOption : 默认命令处理器
type DefaultCMDLineOption struct {
	showVer  bool
	showHelp bool
	//
	logPath  string
	logLevel int
	logSize  int
	//
	coPoolLimt int
	coPoolMax  int
	coPoolMin  int
	//
	virDir string
	//
	configPath string
}

// VarValue 绑定变量
func (cl *DefaultCMDLineOption) VarValue() {
	flag.BoolVar(&cl.showVer, "v", false, "show build version")
	flag.BoolVar(&cl.showHelp, "h", false, "show help")

	flag.StringVar(&cl.logPath, "logPath", "", "log file path")
	flag.IntVar(&cl.logLevel, "logLevel", int(logger.TRACELEVEL), "log level")
	flag.IntVar(&cl.logSize, "logSize", 1024, "log mailbox size")
	flag.StringVar(&cl.virDir, "dir", "", "virtual root directory")

	flag.StringVar(&cl.configPath, "e", "", "config file path")

	flag.IntVar(&cl.coPoolLimt, "colimit", util.MCCOPOOLDEFLIMIT, "maximum stacking limit for cooperative pool tasks")
	flag.IntVar(&cl.coPoolMax, "comax", util.MCCOPOOLDEFMAX, "maximum Limitation of Cooperative Pool")
	flag.IntVar(&cl.coPoolMin, "comin", util.MCCOPOOLDEFMIN, "minimum Limitation of Cooperative Pool")
}

// LineOption : 命令处理函数
func (cl *DefaultCMDLineOption) LineOption() {

	util.PushArgCmd("v", cl.showVer)
	util.PushArgCmd("h", cl.showHelp)
	util.PushArgCmd("logPath", cl.logPath)
	util.PushArgCmd("logLevel", cl.logLevel)
	util.PushArgCmd("logSize", cl.logSize)
	util.PushArgCmd("colimit", cl.coPoolLimt)
	util.PushArgCmd("comax", cl.coPoolMax)
	util.PushArgCmd("comin", cl.coPoolMin)
	util.PushArgCmd("dir", cl.virDir)
	util.PushArgCmd("e", cl.configPath)

	if cl.showVer {
		version.Show()
		os.Exit(0)
	}

	if cl.showHelp {
		cl.usage()
		os.Exit(0)
	}
}

func (cl *DefaultCMDLineOption) usage() {
	f, ferr := os.Open("help.md")
	if ferr != nil {
		panic(fmt.Sprint("error: open help.md file fail:", ferr))
	}
	defer f.Close()
	contents, rerr := ioutil.ReadAll(f)
	if rerr != nil {
		panic(fmt.Sprint("error: open help.md file fail:", rerr))
	}
	fmt.Printf(string(contents))
}
