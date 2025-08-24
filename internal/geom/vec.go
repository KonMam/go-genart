package geom

import "math"

// Vec2 is a 2D vector in logical space [0,1] × [0,1].
type Vec2 struct {
	X, Y float64
}

// Add returns v + u.
func (v Vec2) Add(u Vec2) Vec2 {
	return Vec2{v.X + u.X, v.Y + u.Y}
}

// Sub returns v - u.
func (v Vec2) Sub(u Vec2) Vec2 {
	return Vec2{v.X - u.X, v.Y - u.Y}
}

// Scale returns v * s.
func (v Vec2) Scale(s float64) Vec2 {
	return Vec2{v.X * s, v.Y * s}
}

// Dot product.
func (v Vec2) Dot(u Vec2) float64 {
	return v.X*u.X + v.Y*u.Y
}

// Cross product (scalar in 2D).
func (v Vec2) Cross(u Vec2) float64 {
	return v.X*u.Y - v.Y*u.X
}

// Length returns |v|.
func (v Vec2) Length() float64 {
	return math.Hypot(v.X, v.Y)
}

// Normalize returns v/|v|. Zero vector stays zero.
func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{0, 0}
	}
	return Vec2{v.X / l, v.Y / l}
}
