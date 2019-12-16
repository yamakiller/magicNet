package buffer

import (
	"bytes"

	"github.com/yamakiller/magicLibs/mmath"
)

//NewBuffer doc
//@Summary new a NetBuffer object
//@Method NewBuffer
//@Param int NetBuffer cap
//@Return *NetBuffer
func NewBuffer(cap int) *NetBuffer {
	cap = mmath.Aligned(cap)

	r := &NetBuffer{bytes.NewBuffer([]byte{})}
	if cap > 0 {
		r._data.Grow(cap)
	}
	return r
}

//NetBuffer doc
//@Summary base buffer object
//@Struct NetBuffer
//@Member *bytes.Buffer
type NetBuffer struct {
	_data *bytes.Buffer
}

//Cap doc
//@Summary buffer cap
//@Method Cap
//@Return int
func (slf *NetBuffer) Cap() int {
	return slf._data.Cap()
}

//Len doc
//@Summary buffer Len
//@Method Len
//@Return int
func (slf *NetBuffer) Len() int {
	return slf._data.Len()
}

//Clear doc
//@Summary buffer Clear
//@Method Clear
func (slf *NetBuffer) Clear() {
	slf._data.Reset()
}

//Write doc
//@Summary buffer Write
//@Method Write
//@Param  []byte
//@Return int
//@Return error
func (slf *NetBuffer) Write(d []byte) (int, error) {
	return slf._data.Write(d)
}

//Truncated doc
//@Summary delete buffer n bytes
//@Param truncated data of n size
func (slf *NetBuffer) Truncated(n int) {
	slf._data.Next(n)
}

//Bytes doc
//@Summary Return all bytes
func (slf *NetBuffer) Bytes() []byte {
	return slf._data.Bytes()
}

//Read doc
//@Summary Return n bytes
func (slf *NetBuffer) Read(n int) []byte {
	return slf._data.Next(n)
}
