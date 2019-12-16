package buffer

import "github.com/yamakiller/magicLibs/mmath"

//NewRing doc
//@Summary New a ring buffer
//@Param cap buffer max limit
//@Return RingBuffer
func NewRing(cap int) *NetRingBuffer {
	cap = mmath.Aligned(cap)
	return &NetRingBuffer{_data: make([]byte, cap), _cap: uint(cap)}
}

//NetRingBuffer doc
//@Summary Ring buffer
//@Member []byte data bytes
//@Member  uint ring buffer max limit
//@Member  uint input data size
//@Member  uint output data size
type NetRingBuffer struct {
	_data []byte
	_cap  uint
	_in   uint
	_out  uint
}

//Cap doc
//@Summary Returns buffer max limit
//@Return int
func (slf *NetRingBuffer) Cap() int {
	return int(slf._cap)
}

//Len doc
//@Summary Returns buffer data length
//@Return int
func (slf *NetRingBuffer) Len() int {
	return int(slf._in - slf._out)
}

//Clear doc
//@Summary Clear buffer
func (slf *NetRingBuffer) Clear() {
	slf._in = 0
	slf._out = 0
}

//Bytes doc
//@Summary Return buffer all data
//@Return []byte
func (slf *NetRingBuffer) Bytes() []byte {
	return slf._data
}

//Write doc
//@Summary Writed data to buffer
//@Param []byte data
//@Return writed length
//@Return error
func (slf *NetRingBuffer) Write(p []byte) (n int, err error) {
	length := uint(len(p))
	length = mmath.Min(length, slf._cap-slf._in+slf._out)
	l := mmath.Min(length, slf._cap-(slf._in&(slf._cap-1)))
	copy(slf._data[slf._in&(slf._cap-1):], p[:l])
	copy(slf._data[:length-l], p[l:])
	slf._in += length
	return int(length), nil
}

//Read doc
//@Summary Read data of n size
//@Param  n int readed size
//@Return []byte
func (slf *NetRingBuffer) Read(n int) []byte {
	length := mmath.Min(uint(n), slf._in-slf._out)
	l := mmath.Min(length, slf._cap-(slf._out&(slf._cap-1)))
	result := slf._data[slf._out&(slf._cap-1) : slf._out&(slf._cap-1)+l]
	result = append(result, slf._data[:length-l]...)
	slf._out += length

	return result
}

//Truncated doc
//@Summary Truncated data of n size
//@Para n int size
func (slf *NetRingBuffer) Truncated(n int) {
	length := mmath.Min(uint(n), slf._in-slf._out)
	slf._out += length
}
