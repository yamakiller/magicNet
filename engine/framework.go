package engine

import (
  "magicNet/logger"
)

type framework struct {

}

func (fr *framework) Start() int {
  logger.InitLogger()
  return 0
}

func (fr *framework) Loop() {
}

func (fr *framework) Shutdown() {

}
