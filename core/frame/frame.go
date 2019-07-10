package frame

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月09日 14:36:58
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月10日 16:27:34
 * @Description: 主进程框架基类
 */

/* LineOption ->
Start ->
	LoadEnv ->
		InitService ->
					Wait ->
		CloseService ->
	UnLoadEnv ->
Shutdown */

type startPart interface {
	Start() bool
	Shutdown()
}

type commandLinePart interface {
	LineOption()
}

type envPart interface {
	LoadEnv() bool
	UnLoadEnv()
}

type servicePart interface {
	InitService() bool
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
