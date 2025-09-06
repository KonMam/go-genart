package noise

import "github.com/aquilax/go-perlin"

// PerlinField wraps go-perlin as a ScalarField2D.
type PerlinField struct {
	noise *perlin.Perlin
	scale float64
}

// NewPerlinField creates a new Perlin 2D field.
//
// seed  → determinism
// scale → controls "zoom" (smaller = zoom in, larger = zoom out)
func NewPerlinField(seed int64, scale float64) *PerlinField {
	// alpha, beta, n are tuning params for go-perlin
	// alpha=2, beta=2, n=3 are common defaults
	return &PerlinField{
		noise: perlin.NewPerlin(2, 2, 3, seed),
		scale: scale,
	}
}

// At returns a noise value at (x,y), scaled to [-1,1].
func (p *PerlinField) At(x, y float64) float64 {
	if p.scale <= 0 {
		return 0
	}
	return p.noise.Noise2D(x/p.scale, y/p.scale)
}
