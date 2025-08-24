package noise

import "github.com/ojrac/opensimplex-go"

// SimplexField wraps opensimplex noise as a ScalarField2D.
type SimplexField struct {
	noise opensimplex.Noise // value type
	scale float64
}

// NewSimplexField creates a new Simplex 2D field.
//
// seed  → determinism (use sub-seed for reproducibility)
// scale → controls "zoom" (smaller = zoom in, larger = zoom out)
func NewSimplexField(seed int64, scale float64) *SimplexField {
	return &SimplexField{
		noise: opensimplex.New(seed),
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

// --- 3D variant ---

type SimplexField3D struct {
	noise opensimplex.Noise
	scale float64
}

func NewSimplexField3D(seed int64, scale float64) *SimplexField3D {
	return &SimplexField3D{
		noise: opensimplex.New(seed),
		scale: scale,
	}
}

func (s *SimplexField3D) At(x, y, z float64) float64 {
	if s.scale <= 0 {
		return 0
	}
	return s.noise.Eval3(x/s.scale, y/s.scale, z/s.scale)
}
