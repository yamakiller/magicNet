package core

import (
	"errors"

	"github.com/yamakiller/magicNet/util"
)

// DefaultEnv : 默认的环境变量管理器
type DefaultEnv struct {
}

// LoadEnv : 载入环境变量
func (env *DefaultEnv) LoadEnv() error {
	configPath := util.GetArgString("e", "./env/magicnet.env")
	if configPath == "" {
		return errors.New("enter the environment variable file path  -e <filePath>")
	}

	if err := util.LoadEnv(configPath); err != nil {
		return err
	}

	return nil
}

// UnLoadEnv : 卸载环境变量信息
func (env *DefaultEnv) UnLoadEnv() {
	util.UnLoadEnv()
}
