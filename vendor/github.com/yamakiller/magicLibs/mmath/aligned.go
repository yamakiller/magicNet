package mmath

//Aligned doc
//@Summary aligned byte
func Aligned(n int) int {

	if (n & (n - 1)) != 0 {
		n = n | (n >> 1)
		n = n | (n >> 2)
		n = n | (n >> 4)
		n = n | (n >> 8)
		n = n | (n >> 16)
		n++
	}

	return n
}

//IsPower doc
//@Summary 检测数据是否是２的幂
func IsPower(n int) bool {
	if n < 2 {
		return false
	}

	if (n & (n - 1)) == 0 {
		return true
	}
	return false
}

//Align 对齐字节数
func Align(n uint32, align uint32) uint32 {
	return (n + align - 1) & (^(align - 1))
}
