package util

// IsPower : 判断数值是否椒2的幂
func IsPower(n int) bool {
	if n < 2 {
		return false
	}

	if (n&n - 1) == 0 {
		return true
	}
	return false
}
