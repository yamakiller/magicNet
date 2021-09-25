package util

//NewDFA create an sensitive map
func NewDFA(words []string) *SensitiveMap {
	s := &SensitiveMap{make(map[string]interface{}), false}
	for _, ws := range words {
		sMapTmp := s
		w := []rune(ws)
		wsLen := len(w)
		for i := 0; i < wsLen; i++ {
			t := string(w[i])
			isEnd := false
			if i == (wsLen - 1) {
				isEnd = true
			}
			func(tx string) {
				if _, ok := sMapTmp._sensitiveNode[tx]; !ok {
					sMapTemp := new(SensitiveMap)
					sMapTemp._sensitiveNode = make(map[string]interface{})
					sMapTemp._isend = isEnd
					sMapTmp._sensitiveNode[tx] = sMapTemp
				}
			}(t)
		}
	}
	return s
}

//Target Find Sensitive group
type Target struct {
	Indexes []int
	Len     int
}

//SensitiveMap filer/tire
type SensitiveMap struct {
	_sensitiveNode map[string]interface{}
	_isend         bool
}

//FindOne find first sensitive
func (slf *SensitiveMap) FindOne(text string) (string, bool) {
	content := []rune(text)
	contentLength := len(content)
	result := false
	ta := ""
	for index := range content {
		sMapTmp := slf
		target := ""
		in := index
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp._sensitiveNode[wo]; ok {
				if sMapTmp._sensitiveNode[wo].(*SensitiveMap)._isend {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp._sensitiveNode[wo].(*SensitiveMap) //进入下一层级
				in++
			} else {
				break
			}
		}
		if result {
			ta = target
			break
		}
	}
	return ta, result
}

//FindAll Return target
func (slf *SensitiveMap) FindAll(text string) map[string]*Target {
	content := []rune(text)
	contentLength := len(content)
	result := false

	ta := make(map[string]*Target)
	for index := range content {
		sMapTmp := slf
		target := ""
		in := index
		result = false
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp._sensitiveNode[wo]; ok {
				if sMapTmp._sensitiveNode[wo].(*SensitiveMap)._isend {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp._sensitiveNode[wo].(*SensitiveMap) //进入下一层级
				in++
			} else {
				break
			}
		}
		if result {
			if _, targetInTa := ta[target]; targetInTa {
				ta[target].Indexes = append(ta[target].Indexes, index)
			} else {
				ta[target] = &Target{
					Indexes: []int{index},
					Len:     len([]rune(target)),
				}
			}
		}
	}
	return ta
}
