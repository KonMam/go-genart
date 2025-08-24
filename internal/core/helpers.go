package core

// NewStroke creates a Stroke item with given points and style.
// If closed = true, the path is treated as a closed polygon.
func NewStroke(points []Vec2, closed bool, width float64, color RGBA, alpha float64) Stroke {
	return Stroke{
		Path: Path{
			Points: points,
			Closed: closed,
		},
		Width: width,
		Color: color,
		Alpha: alpha,
	}
}

// NewFill creates a Fill item from a closed polygon.
// Convenience wrapper to avoid repeating boilerplate.
func NewFill(points []Vec2, color RGBA, alpha float64) Fill {
	return Fill{
		Polygon: Path{
			Points: points,
			Closed: true,
		},
		Color: color,
		Alpha: alpha,
	}
}

// Scene helpers

// AddStroke appends a Stroke to a Scene.
func (s *Scene) AddStroke(points []Vec2, closed bool, width float64, color RGBA, alpha float64) {
	s.Items = append(s.Items, NewStroke(points, closed, width, color, alpha))
}

// AddFill appends a Fill to a Scene.
func (s *Scene) AddFill(points []Vec2, color RGBA, alpha float64) {
	s.Items = append(s.Items, NewFill(points, color, alpha))
}
