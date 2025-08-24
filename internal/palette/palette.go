package palette

import (
	"math/rand"

	"genart/internal/core"
)

// Palette is just a slice of RGBA colors.
type Palette []core.RGBA

// Predefined palettes
var (
	Mono = Palette{
		{R: 0, G: 0, B: 0, A: 1},
	}

	Warm = Palette{
		{R: 0.9, G: 0.4, B: 0.2, A: 1},
		{R: 0.95, G: 0.7, B: 0.2, A: 1},
		{R: 0.7, G: 0.2, B: 0.2, A: 1},
	}

	Cool = Palette{
		{R: 0.2, G: 0.5, B: 0.9, A: 1},
		{R: 0.3, G: 0.8, B: 0.7, A: 1},
		{R: 0.1, G: 0.2, B: 0.6, A: 1},
	}

	Rainbow = Palette{
		{R: 1, G: 0, B: 0, A: 1},
		{R: 1, G: 0.5, B: 0, A: 1},
		{R: 1, G: 1, B: 0, A: 1},
		{R: 0, G: 1, B: 0, A: 1},
		{R: 0, G: 0, B: 1, A: 1},
		{R: 0.29, G: 0, B: 0.51, A: 1}, // indigo
		{R: 0.56, G: 0, B: 1, A: 1},    // violet
	}
)

// Pick returns one random color from the palette.
func (p Palette) Pick(rng *rand.Rand) core.RGBA {
	return p[rng.Intn(len(p))]
}

// ByIndex wraps around to allow arbitrary indices.
func (p Palette) ByIndex(i int) core.RGBA {
	return p[i%len(p)]
}
