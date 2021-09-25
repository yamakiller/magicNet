package util

import "unsafe"

//IsLittleEndian doc
//@Method IsLittleEndian @Summary Determine if the system is a small endian
//@Return (bool)
func IsLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return (b == 0x04)
}
