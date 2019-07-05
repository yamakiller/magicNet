package actor

type Process interface {
  SendUsrMessage(pid *PID, message interface{})
  SendSysMessage(pid *PID, message interface{})
  Stop(pid *PID)
}
