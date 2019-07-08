package actor

const pidSetSliceLen = 16

// PIDSet : pid 集对象
type PIDSet struct {
	s []uint32
	m map[uint32]struct{}
}

// NewPIDSet : 新创建一个PID 集
func NewPIDSet(pids ...*PID) *PIDSet {
	var s PIDSet
	for _, pid := range pids {
		s.Add(pid)
	}
	return &s
}

// indexOf 获取某PID在集中的索引
func (p *PIDSet) indexOf(v *PID) int {
	id := v.ID
	for i, pid := range p.s {
		if id == pid {
			return i
		}
	}
	return -1
}

func (p *PIDSet) migrate() {
	p.m = make(map[uint32]struct{}, pidSetSliceLen)
	for _, v := range p.s {
		p.m[v] = struct{}{}
	}
	p.s = p.s[:0]
}

// Add ：PID集中添加一个PID
func (p *PIDSet) Add(v *PID) {
	if p.m == nil {
		if p.indexOf(v) > -1 {
			return
		}

		if len(p.s) < pidSetSliceLen {
			if p.s == nil {
				p.s = make([]uint32, 0, pidSetSliceLen)
			}
			p.s = append(p.s, v.ID)
			return
		}
		p.migrate()
	}
	p.m[v.ID] = struct{}{}
}

// Remove ：PID集中移除一个PID
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
	_, ok := p.m[v.ID]
	if !ok {
		return false
	}
	delete(p.m, v.ID)
	return true
}

// Contains ： PID集中是否包含某个PID
func (p *PIDSet) Contains(v *PID) bool {
	if p.m == nil {
		return p.indexOf(v) != -1
	}
	_, ok := p.m[v.ID]
	return ok
}

// Len ： PID集的长度
func (p *PIDSet) Len() int {
	if p.m == nil {
		return len(p.s)
	}
	return len(p.m)
}

// Clear ：清除PID集
func (p *PIDSet) Clear() {
	if p.m == nil {
		p.s = p.s[:0]
	} else {
		p.m = nil
	}
}

/*Values : 获取整个PID集的Value 并返回数组*/
func (p *PIDSet) Values() []PID {
	if p.Len() == 0 {
		return nil
	}

	r := make([]PID, p.Len())
	if p.m == nil {
		for i, v := range p.s {
			r[i].ID = v
		}
	} else {
		i := 0
		for v := range p.m {
			r[i].ID = v
			i++
		}
	}
	return r
}

// ForEach ：遍历PID集体
func (p *PIDSet) ForEach(f func(i int, pid PID)) {
	var pid PID
	if p.m == nil {
		for i, v := range p.s {
			pid.ID = v
			f(i, pid)
		}
	} else {
		i := 0
		for v := range p.m {
			pid.ID = v
			f(i, pid)
			i++
		}
	}
}

// Clone ：克隆一个PID集
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
