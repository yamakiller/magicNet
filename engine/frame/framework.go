package frame

/* Start -> LineOption
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
	LineOption() int
}

type envPart interface {
	LoadEnv() int
	UnLoadEnv()
}

type servicePart interface {
	InitService()
	CloseService()
}

type waitPart interface {
	Wait()
}

// Framework 主框架接口
type Framework interface {
	startPart
	envPart
	waitPart
	servicePart
}
