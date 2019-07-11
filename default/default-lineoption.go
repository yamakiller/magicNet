package frame

import (
	"flag"
	"fmt"
	"io/ioutil"
	"magicNet/core/version"
	"os"
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
}

// LineOption : 命令处理函数
func (cmdline *DefaultCMDLineOption) LineOption() {
	showVer := flag.Bool("v", false, "show build version")
	if *showVer {
		version.Show()
		os.Exit(0)
	}

	showHelp := flag.Bool("h", false, "show help")
	if *showHelp {
		cmdline.usage()
		os.Exit(0)
	}
}

func (cmdline *DefaultCMDLineOption) usage() {
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
