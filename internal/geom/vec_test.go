package geom

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestVec2Ops(t *testing.T) {
	v := Vec2{3, 4}
	if !almostEqual(v.Length(), 5) {
		t.Errorf("expected length 5, got %f", v.Length())
	}
	n := v.Normalize()
	if !almostEqual(n.Length(), 1) {
		t.Errorf("normalize not unit: %f", n.Length())
	}
}

func TestPolarCartesian(t *testing.T) {
	v := PolarToCartesian(1, math.Pi/2)
	if !almostEqual(v.X, 0) || !almostEqual(v.Y, 1) {
		t.Errorf("expected (0,1), got %v", v)
	}
	r, theta := CartesianToPolar(Vec2{0, 2})
	if !almostEqual(r, 2) || !almostEqual(theta, math.Pi/2) {
		t.Errorf("expected r=2,theta=pi/2, got %f,%f", r, theta)
	}
}

func TestArcEllipse(t *testing.T) {
	arc := Arc(0, 0, 1, 0, math.Pi/2, 5)
	if len(arc) != 5 {
		t.Errorf("expected 5 points, got %d", len(arc))
	}
	ellipse := Ellipse(0, 0, 2, 1, 100)
	if len(ellipse) != 100 {
		t.Errorf("expected 100 points, got %d", len(ellipse))
	}
}
