package aoi

// Point : ---
type Point struct {
	X float32
	Y float32
}

// IsZero : ---
func (pt *Point) IsZero() bool {
	return pt.X == 0 && pt.Y == 0
}

type Size Point
