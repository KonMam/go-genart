package geom

import (
	"math"
	"testing"

	"genart/internal/core"
)

func inUnitBox(p core.Vec2) bool {
	return p.X >= 0 && p.X <= 1 && p.Y >= 0 && p.Y <= 1
}

func TestSquare(t *testing.T) {
	path := Square(0.5, 0.5, 0.6)

	if !path.Closed {
		t.Errorf("square path should be closed")
	}
	if len(path.Points) != 4 {
		t.Errorf("square should have 4 points, got %d", len(path.Points))
	}
	for i, p := range path.Points {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) {
			t.Errorf("point %d is NaN: %+v", i, p)
		}
		if !inUnitBox(p) {
			t.Errorf("point %d out of unit box: %+v", i, p)
		}
	}
}

func TestCircle(t *testing.T) {
	segments := 128
	path := Circle(0.5, 0.5, 0.3, segments)

	if !path.Closed {
		t.Errorf("circle path should be closed")
	}
	if len(path.Points) != segments {
		t.Errorf("circle should have %d points, got %d", segments, len(path.Points))
	}
	for i, p := range path.Points {
		if math.IsNaN(p.X) || math.IsNaN(p.Y) {
			t.Errorf("point %d is NaN: %+v", i, p)
		}
		if !inUnitBox(p) {
			t.Errorf("point %d out of unit box: %+v", i, p)
		}
	}
}
