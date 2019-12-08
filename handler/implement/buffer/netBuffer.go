package buffer

import "bytes"

//New doc
//@Summary new a NetBuffer object
//@Method New
//@Param int NetBuffer cap
//@Return *NetBuffer
func New(cap int) *NetBuffer {
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

func (slf *NetBuffer) Trun(n int) {
	slf._data.Truncate(n)
}

func (slf *NetBuffer) Bytes() []byte {
	return slf._data.Bytes()
}

func (slf *NetBuffer) Read(n int) []byte {
	return slf._data.Next(n)
}
