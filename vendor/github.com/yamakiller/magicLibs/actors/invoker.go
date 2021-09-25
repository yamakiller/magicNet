package actors

type invoker interface {
	invokeSysMessage(msg interface{})
	invokeUsrMessage(msg interface{})
	escalateFailure(reason interface{}, message interface{})
}
