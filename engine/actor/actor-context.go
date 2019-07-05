package actor

type contextState int32

const (
  stateNone contextState = iota
  stateAlive
  stateRestarting
  stateStopping
  stateStopped
)
