package randutil

import (
	"math/rand"
	"testing"
)

func TestRandomRange(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		v := RandomRangeFloat64(rng, 5, 10)
		if v < 5 || v >= 10 {
			t.Fatalf("value %f out of range [5,10)", v)
		}
	}
}

func TestUniformPoint(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	p := UniformPoint(rng)
	if p.X < 0 || p.X > 1 || p.Y < 0 || p.Y > 1 {
		t.Fatalf("point %v not in [0,1]^2", p)
	}
}

func TestJitteredGrid(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	points := JitteredGrid(rng, 5, 0.5)
	if len(points) != 25 {
		t.Fatalf("expected 25 points, got %d", len(points))
	}
	for _, p := range points {
		if p.X < 0 || p.X > 1 || p.Y < 0 || p.Y > 1 {
			t.Fatalf("point %v not in [0,1]^2", p)
		}
	}
}
