package aoi

/*
----------------------------
|  leftTop    |  rightTop  |
|--------------------------
|  leftDonw  |  rightDown |
---------------------------
*/
/*const (
	unknowQuadrant = 0
	rightTop       = 1
	leftTop        = 2
	leftDown       = 3
	rightDonw      = 4
)

// Quadrant : 四叉树 象限信息
type Quadrant int

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
func GenRect(left AFloat, right AFloat, bottom AFloat, top AFloat) Rect {
	return Rect{left: left, right: right, top: top, bottom: bottom,
		midX: (left + (right-left)/2), midY: (bottom + (top-bottom)/2)}
}

// Rect : 范围计算
type Rect struct {
	left   AFloat
	right  AFloat
	top    AFloat
	bottom AFloat
	midX   AFloat
	midY   AFloat
}

// Left :
func (rt *Rect) Left() AFloat {
	return rt.left
}

// Right :
func (rt *Rect) Right() AFloat {
	return rt.right
}

// Bottom :
func (rt *Rect) Bottom() AFloat {
	return rt.bottom
}

// Top :
func (rt *Rect) Top() AFloat {
	return rt.top
}

// MidX :
func (rt *Rect) MidX() AFloat {
	return rt.midX
}

// MidY :
func (rt *Rect) MidY() AFloat {
	return rt.midY
}

// Reset :
func (rt *Rect) Reset() {
	rt.left = 0
	rt.right = 0
	rt.top = 0
	rt.bottom = 0
	rt.midX = 0
	rt.midY = 0
}

// Set :
func (rt *Rect) Set(left AFloat,
	right AFloat,
	bottom AFloat,
	top AFloat) {
	rt.left = left
	rt.right = right
	rt.top = top
	rt.bottom = bottom
	rt.midX = left + (right-left)/2
	rt.midY = bottom - (top-bottom)/2
}

// Contains :
func (rt *Rect) Contains(rect Rect) bool {
	return (rt.left <= rect.left &&
		rt.bottom <= rect.bottom &&
		rect.right <= rt.right &&
		rect.top <= rt.top)
}

// ContainsXY :
func (rt *Rect) ContainsXY(x AFloat, y AFloat) bool {
	return (x >= rt.left && x <= rt.right &&
		y >= rt.bottom && y <= rt.top)
}

// ContainsPoint :
func (rt *Rect) ContainsPoint(point Point) bool {
	return rt.ContainsXY(point.X, point.Y)
}

// Intersects :
func (rt *Rect) Intersects(rect Rect) bool {
	return !(rt.right < rect.left ||
		rect.right < rt.left ||
		rt.top < rect.bottom ||
		rect.top < rt.bottom)
}

// GetQuadrantPoint : 获取四叉树所在象限
func (rt *Rect) GetQuadrantPoint(point Point) Quadrant {
	if rt.ContainsPoint(point) {
		return rt.GetQuadrantPoint2(point)
	}

	return unknowQuadrant
}

// GetQuadrantPoint2 : 获取四叉树象限信息
func (rt *Rect) GetQuadrantPoint2(point Point) Quadrant {
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
}*/
