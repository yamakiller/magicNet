package net

type INetBuffer interface {
	Cap() int
	Len() int
	Clear()
	Trun(n int)
	Bytes() []byte
	Write([]byte) (int, error)
	Read(n int) []byte
}
