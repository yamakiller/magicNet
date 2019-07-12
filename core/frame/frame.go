package frame

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月09日 14:36:58
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月11日 18:25:03
 * @Description: 主进程框架基类
 */

/* LineOption ->
Init ->
	LoadEnv ->
		InitService ->
					Wait ->
		CloseService ->
	UnLoadEnv ->
Destory */

type startPart interface {
	Init() error
	Destory()
}

type commandLinePart interface {
	LineOption()
}

type envPart interface {
	LoadEnv() error
	UnLoadEnv()
}

type servicePart interface {
	InitService() error
	CloseService()
}

type waitPart interface {
	Wait() int
}

// Framework 主框架接口
type Framework interface {
	startPart
	envPart
	waitPart
	servicePart
}

// MakeFrame : 框架制造函数
type MakeFrame func() Framework