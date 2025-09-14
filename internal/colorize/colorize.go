package colorize

import (
	"genart/internal/core"
	"genart/internal/noise"
)

// PickColorFromNoise picks a color from a palette based on a noise value.
func PickColorFromNoise(colors []core.RGBA, noiseField noise.ScalarField2D, x, y, factor float64) core.RGBA {
	if len(colors) == 0 {
		return core.RGBA{R: 0, G: 0, B: 0, A: 1} // fallback to black
	}
	nval := noiseField.At(x*factor, y*factor) // -1..1
	t := (nval + 1) / 2                        // normalize 0..1
	idx := int(t * float64(len(colors)-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(colors) {
		idx = len(colors) - 1
	}
	return colors[idx]
}
