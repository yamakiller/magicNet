package util

import (
	"github.com/yamakiller/magicNet/pool"
)

var (
	defaultCOPool *pool.CoroutinePool
)

// InitCoPool : 初始化协程池
func InitCoPool(limit, max, min int) bool {
	defaultCOPool = &pool.CoroutinePool{TaskLimit: limit, MaxNum: max, MinNum: min}
	defaultCOPool.Start()
	return true
}

// DestoryCoPool : 销毁协咸
func DestoryCoPool() {
	if defaultCOPool != nil {
		defaultCOPool.StopPool()
		defaultCOPool = nil
	}
}

// Go ：调用协程执行协程序
func Go(f func(params []interface{}), params ...interface{}) error {
	if defaultCOPool != nil {
		defaultCOPool.Go(f, params)
	} else {
		go f(params)
	}

	return nil
}
