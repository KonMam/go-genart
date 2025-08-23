package geom

import (
	"math"

	"genart/internal/core"
)

// Square returns a closed path for a square centered at (cx,cy)
func Square(cx, cy, size float64) core.Path {
	hs := size * 0.5
	pts := []core.Vec2{
		{X: cx - hs, Y: cy - hs},
		{X: cx + hs, Y: cy - hs},
		{X: cx + hs, Y: cy + hs},
		{X: cx - hs, Y: cy + hs},
	}
	return core.Path{Points: pts, Closed: true}
}

// Circle returns a closed path for a circle centered at (cx,cy)
func Circle(cx, cy, radius float64, segments int) core.Path {
	if segments < 64 {
		segments = 64
	}
	pts := make([]core.Vec2, 0, segments)
	tau := 2 * math.Pi
	for i := range segments {
		a := float64(i) / float64(segments) * tau
		pts = append(pts, core.Vec2{
			X: cx + radius*math.Cos(a),
			Y: cy + radius*math.Sin(a),
		})
	}
	return core.Path{Points: pts, Closed: true}
}
