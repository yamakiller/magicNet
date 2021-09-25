package mmath

//NewVector3 doc
//@Method NewVector3 @Summary create vector3 object
//@Param  (FValue) x
//@Param  (FValue) y
//@Param  (FValue) z
//@Return (*Vector3)
func NewVector3(x, y, z FValue) *Vector3 {
	return &Vector3{_x: x, _y: y, _z: z}
}

//Vector3 doc
//@Struct Vector3 doc
//@Member (FValue) x
//@Member (FValue) y
//@Member (FValue) z
type Vector3 struct {
	_x FValue
	_y FValue
	_z FValue
}

//Initial doc
//@Summary initialization vector3
//@Param (FValue) x
//@Param (FValue) y
//@Param (FValue) z
func (slf *Vector3) Initial(x, y, z FValue) {
	slf._x = x
	slf._y = y
	slf._z = z
}

//GetX doc
//@Summary return x
//@Return (FValue) x
func (slf *Vector3) GetX() FValue {
	return slf._x
}

//GetY doc
//@Summary return y
//@Return (FValue) y
func (slf *Vector3) GetY() FValue {
	return slf._y
}

//GetZ doc
//@Summary return z
//@Return (FValue)
func (slf *Vector3) GetZ() FValue {
	return slf._z
}

//SetX doc
//@Summary Setting x
//@Param (FValue) x
func (slf *Vector3) SetX(x FValue) {
	slf._x = x
}

//SetY doc
//@Summary Setting y
//@Param (FValue) y
func (slf *Vector3) SetY(y FValue) {
	slf._y = y
}

//SetZ doc
//@Summary Setting z
//@Param (FValue) z
func (slf *Vector3) SetZ(z FValue) {
	slf._z = z
}

func (slf *Vector3) ToVector2() Vector2 {
	return Vector2{_x: slf._x, _y: slf._y}
}
