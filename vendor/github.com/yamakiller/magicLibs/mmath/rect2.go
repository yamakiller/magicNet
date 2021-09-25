package mmath

//NewRect2 new rect 2d
func NewRect2(leftUp, rightDown Vector2) *Rect2 {
	return &Rect2{leftUp, rightDown}
}

//Rect2 2d Rect
type Rect2 struct {
	_pointLeftUp Vector2
	//_pointRightUp   Vector2
	//_pointLeftDown  Vector2
	_pointRightDown Vector2
}

//SetLeftUp Set Left Up point
func (slf *Rect2) SetLeftUp(point Vector2) {
	slf._pointLeftUp = point
}

//GetLeftUp Returns Left up point
func (slf *Rect2) GetLeftUp() *Vector2 {
	return &slf._pointLeftUp
}

/*
//SetRightUp Set Right up point
func (slf *Rect2) SetRightUp(point Vector2) {
	slf._pointRightUp = point
}

//GetRightUp Returns Right up point
func (slf *Rect2) GetRightUp() *Vector2 {
	return &slf._pointRightUp
}*/

/*
//SetLeftDown Set Left down point
func (slf *Rect2) SetLeftDown(point Vector2) {
	slf._pointLeftDown = point
}

//GetLeftDown Returns Left down point
func (slf *Rect2) GetLeftDown() *Vector2 {
	return &slf._pointLeftDown
}*/

//SetRightDown Set Right down point
func (slf *Rect2) SetRightDown(point Vector2) {
	slf._pointRightDown = point
}

//GetRightDown Returns Right down point
func (slf *Rect2) GetRightDown() *Vector2 {
	return &slf._pointRightDown
}

//IsInvolve is in rect
func (slf *Rect2) IsInvolve(point Vector2) bool {
	if point._x >= slf._pointLeftUp._x &&
		point._x < slf._pointRightDown._x &&
		point._y >= slf._pointLeftUp._y &&
		point._y < slf._pointRightDown._y {
		return true
	}
	return false
}

func (slf *Rect2) IsContain(rect Rect2) bool {
	if rect._pointLeftUp.GetX() >= slf._pointLeftUp.GetX() &&
		rect._pointLeftUp.GetY() >= slf._pointLeftUp.GetY() &&
		//rect._pointRightUp.GetX() <= slf._pointRightUp.GetX() &&
		//rect._pointRightUp.GetY() >= slf._pointRightUp.GetY() &&
		//rect._pointLeftDown.GetX() >= slf._pointLeftDown.GetX() &&
		//rect._pointLeftDown.GetY() <= slf._pointLeftDown.GetY() &&
		rect._pointRightDown.GetX() <= slf._pointRightDown.GetX() &&
		rect._pointRightDown.GetY() <= slf._pointRightDown.GetY() {
		return true
	}

	return false
}

func (slf *Rect2) IsIntersect(rect Rect2) bool {
	zx := Abs(slf._pointLeftUp._x + slf._pointRightDown._x - rect._pointLeftUp._x - rect._pointRightDown._x)
	x := Abs(slf._pointLeftUp._x-slf._pointRightDown._x) + Abs(rect._pointLeftUp._x-rect._pointRightDown._x)
	zy := Abs(slf._pointLeftUp._y + slf._pointRightDown._y - rect._pointLeftUp._y - rect._pointRightDown._y)
	y := Abs(slf._pointLeftUp._y-slf._pointRightDown._y) + Abs(rect._pointLeftUp._y-rect._pointRightDown._y)
	if zx <= x && zy <= y {
		return true
	}
	return false
}

/*
//Width Returns Rect width
func (slf *Rect2) Width() FValue {
	return Dist2(slf._pointLeftUp, slf._pointRightUp)
}

//Height Returns Rect height
func (slf *Rect2) Height() FValue {
	return Dist2(slf._pointLeftUp, slf._pointLeftDown)
}*/
