package actors

type dispatcher struct {
	_sch Scheduler
}

func (slf *dispatcher) Schedule(fn func([]interface{})) {
	if slf._sch != nil {
		slf._sch(fn)
	} else {
		go fn(nil)
	}
}
