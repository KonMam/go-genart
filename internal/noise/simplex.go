package noise

import (
	"github.com/ojrac/opensimplex-go"
)

// ScalarField2D is a contract for any scalar-valued 2D field.
type ScalarField2D interface {
	// At returns the field value at (x,y).
	// Typical contract: output in [-1,1].
	At(x, y float64) float64
}

// SimplexField wraps opensimplex noise as a ScalarField2D.
type SimplexField struct {
	noise opensimplex.Noise // store by value, not pointer
	scale float64
}

// NewSimplexField creates a new Simplex noise field.
//
//	seed  → determinism (use a sub-seed for reproducibility)
//	scale → controls "zoom" of the field (smaller = zoom in, larger = zoom out)
func NewSimplexField(seed int64, scale float64) *SimplexField {
	return &SimplexField{
		noise: opensimplex.New(seed), // Noise is a value type
		scale: scale,
	}
}

// At returns a noise value at (x,y), scaled to [-1,1].
func (s *SimplexField) At(x, y float64) float64 {
	if s.scale <= 0 {
		return 0
	}
	return s.noise.Eval2(x/s.scale, y/s.scale)
}
