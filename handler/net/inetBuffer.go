package net

//INetBuffer Network Data buffer interface
type INetBuffer interface {
	Cap() int
	Len() int
	Clear()
	Truncated(n int)
	Bytes() []byte
	Write([]byte) (int, error)
	Read(n int) []byte
}

//INetReceiveBuffer Receive Data buffer interface
type INetReceiveBuffer interface {
	GetBufferCap() int
	GetBufferLen() int
	GetBufferBytes() []byte
	ClearBuffer()
	TrunBuffer(n int)
	WriteBuffer(b []byte) (int, error)
	ReadBuffer(n int) []byte
}
