package main

import (
	"magicNet/engine"
)

func main() {
	fwk := new(engine.Framework)
	if fwk.Start() != 0 {
		fwk.Shutdown()
		return
	}

	fwk.Loop()

	fwk.Shutdown()
}
