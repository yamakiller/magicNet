package aoi

import (
	"unsafe"
)

/*
----------------------------
|  leftTop    |  rightTop  |
|--------------------------
|  leftDonw  |  rightDown |
---------------------------
*/

const (
	unknowQuadrant = 0
	rightTop       = 1
	leftTop        = 2
	leftDown       = 3
	rightDonw      = 4
)

// GenRectPoint : 根据中心点构建区域数据
func GenRectPoint(center Point, half Size) Rect {
	r := Rect{left: center.X - half.X,
		right:  center.X + half.X,
		top:    center.Y + half.Y,
		bottom: center.Y - half.Y}
	r.midX = r.left + (r.right-r.left)/2
	r.midY = r.bottom + (r.top-r.bottom)/2
	return r
}

// GenRect : 构建一个区域
func GenRect(left float32, right float32, bottom float32, top float32) Rect {
	return Rect{left: left, right: right, top: top, bottom: bottom,
		midX: (left + (right-left)/2), midY: (bottom + (top-bottom)/2)}
}

type Rect struct {
	left   float32
	right  float32
	top    float32
	bottom float32
	midX   float32
	midY   float32
}

type Quadrant int

// Left :
func (rt *Rect) Left() float32 {
	return rt.left
}

// Right :
func (rt *Rect) Right() float32 {
	return rt.right
}

// Bottom :
func (rt *Rect) Bottom() float32 {
	return rt.bottom
}

// Top :
func (rt *Rect) Top() float32 {
	return rt.top
}

// MidX :
func (rt *Rect) MidX() float32 {
	return rt.midX
}

// MidY :
func (rt *Rect) MidY() float32 {
	return rt.midY
}

// Reset :
func (rt *Rect) Reset(args ...float32) {

	v := []float32(args)
	if len(v) == 0 {
		rt.left = 0
		rt.right = 0
		rt.top = 0
		rt.bottom = 0
		rt.midX = 0
		rt.midY = 0
	} else {
		left := v[0]
		right := v[1]
		bottom := v[2]
		top := v[3]

		rt.left = left
		rt.right = right
		rt.top = top
		rt.bottom = bottom
		rt.midX = left + (right-left)/2
		rt.midY = bottom - (top-bottom)/2
	}
}

// Contains : --
func (rt *Rect) Contains(in unsafe.Pointer) bool {
	if inRect := (*Rect)(in); inRect != nil {
		return rt.containsRect(inRect)
	}

	if inPoint := (*Point)(in); inPoint != nil {
		return rt.containsPoint(inPoint)
	}

	panic("Contains not Rect or Point")
}

// Intersects : ---
func (rt *Rect) Intersects(rect *Rect) bool {
	return !(rt.right < rect.left ||
		rect.right < rt.left ||
		rt.top < rect.bottom ||
		rect.top < rt.bottom)
}

// GetQuadrant : --
func (rt *Rect) GetQuadrant(point *Point) Quadrant {
	if rt.containsPoint(point) {
		return rt.GetQuadrant2(point)
	}

	return unknowQuadrant
}

// GetQuadrant2 :--
func (rt *Rect) GetQuadrant2(point *Point) Quadrant {
	if point.Y >= rt.midY {
		if point.X >= rt.midX {
			return rightTop
		}
		return leftTop
	}

	if point.X >= rt.midX {
		return rightDonw
	}

	return leftDown
}

func (rt *Rect) containsRect(rect *Rect) bool {
	return (rt.left <= rect.left &&
		rt.bottom <= rect.bottom &&
		rect.right <= rt.right &&
		rect.top <= rt.top)
}

func (rt *Rect) containsPoint(point *Point) bool {
	return rt.containsXY(point.X, point.Y)
}

func (rt *Rect) containsXY(x float32, y float32) bool {
	return (x >= rt.left && x <= rt.right &&
		y >= rt.bottom && y <= rt.top)
}
