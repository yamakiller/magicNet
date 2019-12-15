package buffer

func NewRBuffer(cap uint) *NetRBuffer {
	// if size not pow of 2 round up it
	if (cap & (cap - 1)) != 0 {
		cap = cap | (cap >> 1)
		cap = cap | (cap >> 2)
		cap = cap | (cap >> 4)
		cap = cap | (cap >> 8)
		cap = cap | (cap >> 16)
		cap++
	}

	return &NetRBuffer{_data: make([]byte, cap), _cap: cap}
}

type NetRBuffer struct {
	_data []byte
	_cap  uint
	_in   uint
	_out  uint
}

func (slf *NetRBuffer) Cap() int {
	return int(slf._cap)
}

func (slf *NetRBuffer) Len() int {
	return int(slf._in - slf._out)
}

func (slf *NetRBuffer) Clear() {
	slf._in = 0
	slf._out = 0
}

func (slf *NetRBuffer) Bytes() []byte {
	return slf._data
}

func (slf *NetRBuffer) Write(p []byte) (n int, err error) {
	length := uint(len(p))
	length = min(length, slf._cap-slf._in+slf._out)
	l := min(length, slf._cap-(slf._in&(slf._cap-1)))
	copy(slf._data[slf._in&(slf._cap-1):], p[:l])
	copy(slf._data[:length-l], p[l:])
	slf._in += length
	return int(length), nil
}

func (slf *NetRBuffer) Read(n int) []byte {
	length := min(uint(n), slf._in-slf._out)
	l := min(length, slf._cap-(slf._out&(slf._cap-1)))
	result := slf._data[slf._out&(slf._cap-1) : slf._out&(slf._cap-1)+1]
	result = append(result, slf._data[:length-l]...)
	slf._out += length

	return result
}

func (slf *NetRBuffer) Turn(n int) {
	length := min(uint(n), slf._in-slf._out)
	slf._out += length
}

func min(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}
