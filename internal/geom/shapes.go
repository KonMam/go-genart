package geom

import "math"

// PolarToCartesian converts polar coordinates (r, θ) to a Vec2.
func PolarToCartesian(r, theta float64) Vec2 {
	return Vec2{
		X: r * math.Cos(theta),
		Y: r * math.Sin(theta),
	}
}

// CartesianToPolar converts (x,y) into (r,θ).
func CartesianToPolar(v Vec2) (r, theta float64) {
	r = v.Length()
	theta = math.Atan2(v.Y, v.X)
	return
}

// Arc returns n points approximating an arc of a circle.
// cx,cy = center, r = radius, start,end = angles in radians.
func Arc(cx, cy, r, start, end float64, n int) []Vec2 {
	pts := make([]Vec2, n)
	dt := (end - start) / float64(n-1)
	for i := 0; i < n; i++ {
		theta := start + float64(i)*dt
		pts[i] = Vec2{
			X: cx + r*math.Cos(theta),
			Y: cy + r*math.Sin(theta),
		}
	}
	return pts
}

// Ellipse returns n points approximating an ellipse.
// cx,cy = center, rx,ry = radii.
func Ellipse(cx, cy, rx, ry float64, n int) []Vec2 {
	pts := make([]Vec2, n)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		pts[i] = Vec2{
			X: cx + rx*math.Cos(theta),
			Y: cy + ry*math.Sin(theta),
		}
	}
	return pts
}
