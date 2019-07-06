package actor

const pidSetSliceLen = 16

type PIDSet struct {
  s []uint32
  m map[uint32]struct{}
}

func NewPIDSet(pids ...*PID) *PIDSet {
  var s PIDSet
  for _, pid := range pids {
    s.Add(pid)
  }
  return &s
}

func (p *PIDSet) indexOf(v *PID) int {
  id := v.Id
  for i, pid := range p.s {
    if id == pid {
      return i
    }
  }
  return -1
}

func (p *PIDSet) migrate(){
  p.m = make(map[uint32]struct{}, pidSetSliceLen)
  for _, v := range p.s {
    p.m[v] = struct{}{}
  }
  p.s = p.s[:0]
}

func (p *PIDSet) Add(v *PID) {
  if p.m == nil {
    if p.indexOf(v) > -1 {
      return
    }

    if len(p.s) < pidSetSliceLen {
      if p.s == nil {
        p.s = make([]uint32, 0, pidSetSliceLen)
      }
      p.s = append(p.s, v.Id)
      return
    }
    p.migrate()
  }
  p.m[v.Id] = struct{}{}
}

func (p *PIDSet) Remove(v *PID) bool {
  if p.m == nil {
    i := p.indexOf(v)
    if i == -1 {
      return false
    }
    l := len(p.s) - 1
    p.s[i] = p.s[l]
    p.s = p.s[:l]
    return true
  }
  _, ok := p.m[v.Id]
  if !ok {
    return false
  }
  delete(p.m, v.Id)
  return true
}

func (p *PIDSet) Contains(v *PID) bool {
  if p.m == nil {
    return p.indexOf(v) != -1
  }
  _, ok := p.m[v.Id]
  return ok
}

func (p *PIDSet) Len() int {
  if p.m == nil {
    return len(p.s)
  }
  return len(p.m)
}

func (p *PIDSet) Clear() {
  if p.m == nil {
    p.s = p.s[:0]
  } else {
    p.m = nil
  }
}

func (p *PIDSet)  Values() []PID {
  if p.Len() == 0 {
    return nil
  }

  r := make([]PID, p.Len())
  if p.m == nil {
    for i,v := range p.s {
        r[i].Id = v
    }
  } else {
    i := 0
    for v := range p.m {
      r[i].Id = v
      i++
    }
  }
  return r
}

func (p *PIDSet) ForEach(f func(i int, pid PID)) {
  var pid PID
  if p.m == nil {
    for i, v := range p.s {
      pid.Id = v
      f(i, pid)
    }
  } else {
    i := 0
    for v := range p.m {
      pid.Id = v
      f(i, pid)
      i++
    }
  }
}

func (p *PIDSet) Clone() *PIDSet {
  var s PIDSet
	if p.s != nil {
		s.s = make([]uint32, len(p.s))
		for i, v := range p.s {
			s.s[i] = v
		}
	}
	if p.m != nil {
		s.m = make(map[uint32]struct{}, len(p.m))
		for v := range p.m {
			s.m[v] = struct{}{}
		}
	}
	return &s
}
