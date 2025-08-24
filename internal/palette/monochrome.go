package palette

import "genart/internal/core"

// Monochrome generates a monochromatic palette.
// Keeps hue & saturation fixed, varies lightness.
func Monochrome(base core.RGBA, n int) []core.RGBA {
	if n < 2 {
		n = 2
	}

	h, s, l := RGBToHSL(base.R, base.G, base.B)
	colors := make([]core.RGBA, n)

	for i := 0; i < n; i++ {
		f := float64(i) / float64(n-1)
		// lightness range: darker to lighter around base
		newL := clamp(l*0.3+f*0.7, 0, 1)
		r, g, b := HSLToRGB(h, s, newL)
		colors[i] = core.RGBA{R: r, G: g, B: b, A: 1}
	}
	return colors
}
