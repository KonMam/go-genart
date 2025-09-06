package palette

import (
	"genart/internal/core"
	"math"
)

func Analogous(base core.RGBA, n int) []core.RGBA {
	if n < 3 {
		n = 3
	}

	h, s, l := RGBToHSL(base.R, base.G, base.B)
	colors := make([]core.RGBA, n)

	hues := []float64{math.Mod(h-1.0/12.0, 1.0), h, math.Mod(h+1.0/12.0, 1.0)}

	for i := 0; i < n; i++ {
		// Cycle through the three main hues
		hue := hues[i%3]

		// Vary lightness and saturation
		lightness := clamp(l*0.3+float64(i)/float64(n-1)*0.7, 0, 1)
		saturation := clamp(s*0.5+float64(i)/float64(n-1)*0.5, 0, 1)

		r, g, b := HSLToRGB(hue, saturation, lightness)
		colors[i] = core.RGBA{R: r, G: g, B: b, A: 1}
	}

	return colors
}
