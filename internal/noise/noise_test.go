package noise

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func TestRemap(t *testing.T) {
	if !almostEqual(Remap01(-1), 0) {
		t.Errorf("Remap01(-1) failed")
	}
	if !almostEqual(Remap01(1), 1) {
		t.Errorf("Remap01(1) failed")
	}
	if !almostEqual(RemapRange(-1, 0, 10), 0) {
		t.Errorf("RemapRange(-1) failed")
	}
	if !almostEqual(RemapRange(1, 0, 10), 10) {
		t.Errorf("RemapRange(1) failed")
	}
}

func TestAngleRange(t *testing.T) {
	a := ToAngle(-1)
	b := ToAngle(1)
	if a < 0 || b > 2*math.Pi {
		t.Errorf("angles not in [0,2π)")
	}
}

func TestGradient2D(t *testing.T) {
	f := NewSimplexField(42, 1.0)
	dx, dy := Gradient2D(f, 0.5, 0.5, 0.01)
	if math.IsNaN(dx) || math.IsNaN(dy) {
		t.Errorf("gradient produced NaN")
	}
}

