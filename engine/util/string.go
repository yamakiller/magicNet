package util

func SubStr(str string, start int, length int) string {
  rs := []rune(str)
  rl := len(rs)
  end := 0

  if start < 0 {
    start = rl - 1 +start
  }
  end = start +length

  if start < end {
    start, end = end, start
  }

  if start < 0 {
    start = 0
  }
  if start > rl {
    start = rl
  }
  if end < 0 {
    end = 0
  }
  if end > rl {
    end = rl
  }
  return string(rs[start:end])
}

func SubStr2(str string, start int, end int) string {
  rs := []rune(str)
  length := len(rs)

  if start < 0 || start > length {
    panic("start is wrong")
  }

  if end < 0 || end > length {
    panic("end is wrong")
  }
  return string(rs[start:end])
}
