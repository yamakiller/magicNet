package jslib

import (
	"runtime"

	"github.com/robertkrimen/otto"
)

func jsruntimeBundle(js otto.FunctionCall) otto.Value {
	truntime, _ := js.Otto.Object(`({})`)
	truntime.Set("readMemStats", readMemStats)
	truntime.Set("numGoroutine", numGoroutine)
	truntime.Set("numCPU", numCPU)
	truntime.Set("platform", platform)

	vruntime, err := otto.ToValue(truntime)
	if err != nil {
		panic(err)
	}

	return vruntime
}

func readMemStats(call otto.FunctionCall) otto.Value {
	var mst runtime.MemStats
	runtime.ReadMemStats(&mst)
	result, _ := call.Otto.Object(`({})`)
	result.Set("Sys", mst.Sys)
	result.Set("Lookups", mst.Lookups)
	result.Set("Mallocs", mst.Mallocs)
	result.Set("Frees", mst.Frees)
	result.Set("HeapAlloc", mst.HeapAlloc)
	result.Set("HeapSys", mst.HeapSys)
	result.Set("HeapIdle", mst.HeapIdle)
	result.Set("HeapInuse", mst.HeapInuse)
	result.Set("HeapReleased", mst.HeapReleased)
	result.Set("HeapObjects", mst.HeapObjects)
	result.Set("StackInuse", mst.StackInuse)
	result.Set("StackSys", mst.StackSys)
	result.Set("MSpanInuse", mst.MSpanInuse)
	result.Set("MSpanSys", mst.MSpanSys)
	result.Set("MCacheInuse", mst.MCacheInuse)
	result.Set("MCacheSys", mst.MCacheSys)
	result.Set("BuckHashSys", mst.BuckHashSys)
	result.Set("GCSys", mst.GCSys)
	result.Set("OtherSys", mst.OtherSys)
	result.Set("NextGC", mst.NextGC)
	result.Set("LastGC", mst.LastGC)
	result.Set("PauseTotalNs", mst.PauseTotalNs)
	result.Set("PauseNs", mst.PauseNs)
	result.Set("PauseEnd", mst.PauseEnd)
	result.Set("NumForcedGC", mst.NumForcedGC)
	result.Set("GCCPUFraction", mst.GCCPUFraction)
	result.Set("BySize", mst.BySize)
	result.Set("Alloc", mst.Alloc)
	result.Set("TotalAlloc", mst.TotalAlloc)

	vmst, err := otto.ToValue(result)
	if err != nil {
		panic(err)
	}

	return vmst
}

func numGoroutine(call otto.FunctionCall) otto.Value {
	vnum, err := otto.ToValue(runtime.NumGoroutine())
	if err != nil {
		panic(err)
	}
	return vnum
}

func numCPU(call otto.FunctionCall) otto.Value {
	vnum, err := otto.ToValue(runtime.NumCPU())
	if err != nil {
		panic(err)
	}
	return vnum
}

func platform(call otto.FunctionCall) otto.Value {
	vstr, err := otto.ToValue(runtime.GOOS)
	if err != nil {
		panic(err)
	}
	return vstr
}
