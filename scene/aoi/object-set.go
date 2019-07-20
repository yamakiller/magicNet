package aoi

//? 需要在效率上优化
type objectSet struct {
	num   int
	set   map[uint32]interface{}
	count int
}

func (oset *objectSet) Add(id uint32) {
	if oset.set[id] != nil {
		oset.set[id] = struct{}{}
		oset.count++
	}
}

func (oset *objectSet) Del(id uint32) {
	if oset.set[id] != nil {
		delete(oset.set, id)
		oset.count--
	}
}

func (oset *objectSet) ToArray() []uint32 {
	if oset.count == 0 {
		return nil
	}

	i := 0
	r := make([]uint32, oset.count)
	for k := range oset.set {
		r[i] = k
		i++
	}
	return r
}

func (oset *objectSet) Rest() {
	if oset.count > 100 {
		if oset.num > 0 {
			oset.set = make(map[uint32]interface{}, oset.num)
		} else {
			oset.set = make(map[uint32]interface{})
		}
	} else {
		for k := range oset.set {
			delete(oset.set, k)
		}
	}
	oset.count = 0
}

func (oset *objectSet) ToCopy(dst *objectSet) {
	for k := range oset.set {
		dst.Add(k)
	}
}

func (oset *objectSet) Except(oset2 *objectSet, dst *objectSet) {
	for k := range oset.set {
		if oset2.set[k] == nil {
			dst.Add(k)
		}
	}

	for k := range oset2.set {
		if oset.set[k] == nil {
			dst.Add(k)
		}
	}
}
