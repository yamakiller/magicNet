package util

func IsPower(n int) bool {
  if (n < 2) {
    return false
  }

  if (n & n - 1) == 0 {
    return true
  }
  return false
}
