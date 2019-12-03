package frame

/*
 * @Author: mirliang@my.cn
 * @Date: 2019年07月09日 14:36:58
 * @LastEditors: mirliang@my.cn
 * @LastEditTime: 2019年07月20日 18:15:02
 * @Description: 主进程框架基类
 */

/*
Init ->
		InitService ->
					Wait ->
		CloseService ->
Destory */

type bootPart interface {
	Initial() error
	Destory()
}

type servPart interface {
	InitService() error
	CloseService()
}

type waitPart interface {
	Enter()
	Wait() int
}

//Framework desc:
//@Interface Framework desc: system frame
type Framework interface {
	bootPart
	waitPart
	servPart
}

// SpawnFrame desc:
// @type SpawnFrame desc: create main framework
type SpawnFrame func() Framework
