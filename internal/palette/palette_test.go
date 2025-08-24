package palette

import (
	"math/rand"
	"testing"
)

func TestPick(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	c := Warm.Pick(rng)
	if c.R == 0 && c.G == 0 && c.B == 0 {
		t.Fatalf("expected non-black color from Warm palette")
	}
}

func TestByIndex(t *testing.T) {
	c := Rainbow.ByIndex(100) // should wrap around
	if c.A != 1 {
		t.Fatalf("expected alpha=1, got %f", c.A)
	}
}
