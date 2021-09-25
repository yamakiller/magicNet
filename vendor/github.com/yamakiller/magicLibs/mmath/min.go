package mmath

//Min doc
//@Summary Returns x or y min value
//@Param x uint
//@Param y uint
//@Return uint
func Min(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}
