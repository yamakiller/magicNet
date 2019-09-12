package comparator

import "time"

// Comparator ：
type Comparator func(a, b interface{}) int

//StringComparator default a fast comparison on strings
func StringComparator(a, b interface{}) int {
	s1 := a.(string)
	s2 := b.(string)

	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

// IntComparator default a fast comparison on int
func IntComparator(a, b interface{}) int {
	aAss := a.(int)
	bAss := b.(int)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Int8Comparator default a fast comparison on int8
func Int8Comparator(a, b interface{}) int {
	aAss := a.(int8)
	bAss := b.(int8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Int16Comparator default a fast comparison on int16
func Int16Comparator(a, b interface{}) int {
	aAss := a.(int16)
	bAss := b.(int16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Int32Comparator default a fast comparison on int32
func Int32Comparator(a, b interface{}) int {
	aAss := a.(int32)
	bAss := b.(int32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Int64Comparator default a fast comparison on int64
func Int64Comparator(a, b interface{}) int {
	aAss := a.(int64)
	bAss := b.(int64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// UIntComparator default a fast comparison on uint
func UIntComparator(a, b interface{}) int {
	aAss := a.(uint)
	bAss := b.(uint)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// UInt8Comparator default a fast comparison on uint8
func UInt8Comparator(a, b interface{}) int {
	aAss := a.(uint8)
	bAss := b.(uint8)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// UInt16Comparator default a fast comparison on uint16
func UInt16Comparator(a, b interface{}) int {
	aAss := a.(uint16)
	bAss := b.(uint16)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// UInt32Comparator default a fast comparison on uint32
func UInt32Comparator(a, b interface{}) int {
	aAss := a.(uint32)
	bAss := b.(uint32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// UInt64Comparator default a fast comparison on uint64
func UInt64Comparator(a, b interface{}) int {
	aAss := a.(uint64)
	bAss := b.(uint64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Float32Comparator default a fast comparison on float32
func Float32Comparator(a, b interface{}) int {
	aAss := a.(float32)
	bAss := b.(float32)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// Float64Comparator default a fast comparison on float64
func Float64Comparator(a, b interface{}) int {
	aAss := a.(float64)
	bAss := b.(float64)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// ByteComparator default a fast comparison on byte
func ByteComparator(a, b interface{}) int {
	aAss := a.(byte)
	bAss := b.(byte)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// RuneComparator default a fast comparison on  time.Time
func RuneComparator(a, b interface{}) int {
	aAss := a.(rune)
	bAss := b.(rune)
	switch {
	case aAss > bAss:
		return 1
	case aAss < bAss:
		return -1
	default:
		return 0
	}
}

// TimeComparator default a fast comparison on rune
func TimeComparator(a, b interface{}) int {
	aAss := a.(time.Time)
	bAss := b.(time.Time)
	switch {
	case aAss.After(bAss):
		return 1
	case aAss.Before(bAss):
		return -1
	default:
		return 0
	}
}
