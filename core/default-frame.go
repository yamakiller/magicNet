package core

import (
	"github.com/yamakiller/magicNet/core/frame"
)

// DefaultFrame : 默认框架
type DefaultFrame struct {
	DefaultStart
	DefaultCMDLineOption
	DefaultEnv
	DefaultService
	DefaultLoop
}

var (
	_ frame.Framework = &DefaultFrame{}
)
