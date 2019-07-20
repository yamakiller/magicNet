package aoi

import (
	"math"
)

// Vector2 : 二维顶点
type Vector2 struct {
	X float64
	Y float64
}

// Distance2 : 计算距离
func Distance2(p1 Vector2, p2 Vector2) float64 {
	return math.Pow(((p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y)), 0.5)
}
