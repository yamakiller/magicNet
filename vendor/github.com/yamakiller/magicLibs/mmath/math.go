package mmath

import (
	"unsafe"
)

//Dist2 2d Point x1 => x2 distance
func Dist2(p1 Vector2, p2 Vector2) FValue {
	return InvSqrt((p2._x-p1._x)*(p2._x-p1._x)) + InvSqrt((p2._y-p1._y)*(p2._y-p1._y))
}

//Dist3 3d Point x1 => x2 distance
func Dist3(p1 Vector3, p2 Vector3) FValue {
	return InvSqrt((p2._x-p1._x)*(p2._x-p1._x)) + InvSqrt((p2._y-p1._y)*(p2._y-p1._y)) + InvSqrt((p2._z-p1._z)*(p2._z-p1._z))
}

//InvSqrt Quake3 sqrt
func InvSqrt(x FValue) FValue {
	xhalf := float64(x) * float64(0.5)
	i := (*(*int)(unsafe.Pointer(&x)))
	i = 0x5f3759df - (i >> 1)
	x = (*(*FValue)(unsafe.Pointer(&i)))
	x = FValue(float64(x) * (float64(1.5) - xhalf*float64(x)*float64(x)))
	return x
}

//Abs abs
func Abs(x FValue) FValue {
	if x >= 0 {
		return x
	}

	return -x
}
