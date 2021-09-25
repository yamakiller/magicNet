package mmath

//NewVector2 doc
//@Method NewVector2 @Summary create vector2 object
//@Param  (FValue) x
//@Param  (FValue) y
//@Return (*Vector2)
func NewVector2(x, y FValue) *Vector2 {
	return &Vector2{_x: x, _y: y}
}

//Vector2 doc
//@Struct Vector2 doc
//@Member (float64) x
//@Member (float64) y
type Vector2 struct {
	_x FValue
	_y FValue
}

//Initial doc
//@Method Initial @Summary initialization vector2
//@Param (FValue) x
//@Param (FValue) y
func (slf *Vector2) Initial(x, y FValue) {
	slf._x = x
	slf._y = y
}

//GetX doc
//@Summary return x
//@Return (FValue) x
func (slf *Vector2) GetX() FValue {
	return slf._x
}

//GetY doc
//@Method GetY @Summary return y
//@Return (FValue) y
func (slf *Vector2) GetY() FValue {
	return slf._y
}

//SetX doc
//@Summary Setting x
//@Param (FValue) x
func (slf *Vector2) SetX(x FValue) {
	slf._x = x
}

//SetY doc
//@Summary Setting y
//@Param (FValue) y
func (slf *Vector2) SetY(y FValue) {
	slf._y = y
}

//ToVector3 doc
//@Summary To Vector3
//@Return Vector3
func (slf *Vector2) ToVector3() Vector3 {
	return Vector3{_x: slf._x, _y: slf._y, _z: 0}
}
