package frame

import (
	"flag"
	"fmt"
	"magicNet/engine/util"
	"strings"
)

type defaultEnv struct {
}

func (env *defaultEnv) LoadEnv() bool {
	configPath := ""
	flag.StringVar(&configPath, "e", "./conf/magicnet.conf", "config full path")
	if strings.Compare(configPath, "") == 0 {
		panic("enter the environment variable file path  -e <filePath>")
	}

	if util.LoadEnv(configPath) != 0 {
		panic(fmt.Sprint("failed to open environment variable configuration file:", configPath))
	}

	return true
}

func (env *defaultEnv) UnLoadEnv() {
	util.UnLoadEnv()
}
