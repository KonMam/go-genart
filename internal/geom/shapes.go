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

// Circle is a convenience wrapper for Ellipse with rx = ry.
// Ensures a minimum segment count for smoothness.
func Circle(cx, cy, r float64, n int) []Vec2 {
	if n < 64 {
		n = 64
	}
	return Ellipse(cx, cy, r, r, n)
}

// Polygon returns a regular n-gon centered at (cx,cy) with circumradius r.
func Polygon(cx, cy, r float64, n int) []Vec2 {
	pts := make([]Vec2, n)
	tau := 2 * math.Pi
	for i := 0; i < n; i++ {
		theta := tau * float64(i) / float64(n)
		pts[i] = Vec2{
			X: cx + r*math.Cos(theta),
			Y: cy + r*math.Sin(theta),
		}
	}
	return pts
}

// Square is a convenience wrapper for Polygon with n=4.
func Square(cx, cy, size float64) []Vec2 {
	// size is side length, so circumradius = size/√2
	r := size / math.Sqrt2
	return Polygon(cx, cy, r, 4)
}

// --- Transform helpers ---

// Translate shifts all points by (dx,dy).
func Translate(pts []Vec2, dx, dy float64) []Vec2 {
	out := make([]Vec2, len(pts))
	for i, p := range pts {
		out[i] = Vec2{X: p.X + dx, Y: p.Y + dy}
	}
	return out
}

// Scale scales all points relative to origin (0,0).
func Scale(pts []Vec2, sx, sy float64) []Vec2 {
	out := make([]Vec2, len(pts))
	for i, p := range pts {
		out[i] = Vec2{X: p.X * sx, Y: p.Y * sy}
	}
	return out
}

// Rotate rotates all points around origin (0,0) by angle radians.
func Rotate(pts []Vec2, angle float64) []Vec2 {
	out := make([]Vec2, len(pts))
	c, s := math.Cos(angle), math.Sin(angle)
	for i, p := range pts {
		out[i] = Vec2{
			X: p.X*c - p.Y*s,
			Y: p.X*s + p.Y*c,
		}
	}
	return out
}
