package frame

import (
	"errors"
	"flag"
	"fmt"
	"magicNet/engine/util"
	"strings"
)

// DefaultEnv : 默认的环境变量管理器
type DefaultEnv struct {
}

// LoadEnv : 载入环境变量
func (env *DefaultEnv) LoadEnv() error {
	configPath := ""
	flag.StringVar(&configPath, "e", "./conf/magicnet.conf", "config full path")
	if strings.Compare(configPath, "") == 0 {
		return errors.New("enter the environment variable file path  -e <filePath>")
	}

	if util.LoadEnv(configPath) != 0 {
		return errors.New(fmt.Sprint("failed to open environment variable configuration file:", configPath))
	}

	return nil
}

// UnLoadEnv : 卸载环境变量信息
func (env *DefaultEnv) UnLoadEnv() {
	util.UnLoadEnv()
}
