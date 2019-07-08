package actor

// Process : 处理模块基础接口
type Process interface {
	SendUsrMessage(pid *PID, message interface{})
	SendSysMessage(pid *PID, message interface{})
	OverloadUsrMessage() int
	Stop(pid *PID)
}
