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

type defaultCMDLineOption struct {
}

func (cmdline *defaultCMDLineOption) LineOption() {
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

func (cmdline *defaultCMDLineOption) usage() {
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
