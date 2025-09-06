package palette

import (
	"math"

	"genart/internal/core"
)

// SplitComplementary generates a vivid split-complementary palette.
func SplitComplementary(base core.RGBA, n int) []core.RGBA {
	if n < 3 {
		n = 3
	}

	h, s, _ := RGBToHSL(base.R, base.G, base.B)
	colors := make([]core.RGBA, n)

	// main hue + split complementary (±30° from opposite)
	hues := []float64{
		h,
		math.Mod(h+0.5-1.0/12.0, 1.0),
		math.Mod(h+0.5+1.0/12.0, 1.0),
	}

	for i := 0; i < n; i++ {
		hue := hues[i%3]

		// push saturation high, keep lightness around mid
		saturation := clamp(s*0.9+(float64(i)/float64(n-1))*0.1, 0.6, 1.0)
		lightness := clamp(0.4+(float64(i)/float64(n-1))*0.3, 0, 1)

		r, g, b := HSLToRGB(hue, saturation, lightness)
		colors[i] = core.RGBA{R: r, G: g, B: b, A: 1}
	}

	return colors
}
