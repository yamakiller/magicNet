package net

type INetBuffer interface {
	Cap() int
	Len() int
	Clear()
	Write([]byte) (int, error)
}
