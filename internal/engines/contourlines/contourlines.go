package contourlines

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/noise"
)

type Engine struct{}

func (Engine) Name() string { return "contourlines" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	lines := int(pick(params, "lines", 3000))
	steps := int(pick(params, "steps", 500))
	scale := pick(params, "scale", 0.01)
	step := pick(params, "step", 0.0008)
	resetProb := pick(params, "resetProb", 0.005)
	dotSize := pick(params, "dotSize", 0.0015)
	// paletteID := int(pick(params, "palette", 1)) // default warm

	field := noise.NewSimplexField(rng.Int63(), scale)
	scene := core.Scene{}

	// colors := selectPalette(paletteID)

	for i := 0; i < lines; i++ {
		x, y := rng.Float64(), rng.Float64()
		thetaPrev := rng.Float64() * 2 * math.Pi

		for j := 0; j < steps; j++ {
			if rng.Float64() < resetProb {
				break
			}

			thetaNoise := field.At(x/100, y/100) * math.Pi * 2
			theta := 0.9*thetaPrev + 0.1*thetaNoise
			thetaPrev = theta

			x += math.Cos(theta) * step
			y += math.Sin(theta) * step

			if x < 0 || x > 1 || y < 0 || y > 1 {
				break
			}

			// Drop a dot at each step
			dot := circleAt(x, y, dotSize, 12)
			scene.Items = append(scene.Items, core.Fill{
				Polygon: dot,
				Color:   colors[rng.Intn(len(colors))],
				Alpha:   0.5 + rng.Float64()*0.4, // 0.5–0.9
			})
		}
	}

	return scene, nil
}

// --- helpers ---

func circleAt(cx, cy, r float64, segs int) core.Path {
	pts := make([]core.Vec2, segs)
	for i := 0; i < segs; i++ {
		a := 2 * math.Pi * float64(i) / float64(segs)
		pts[i] = core.Vec2{
			X: cx + r*math.Cos(a),
			Y: cy + r*math.Sin(a),
		}
	}
	return core.Path{Points: pts, Closed: true}
}

func pick(m map[string]float64, k string, def float64) float64 {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}

func selectPalette(id int) []core.RGBA {
	switch id {
	case 1: // warm
		return []core.RGBA{
			{R: 0.9, G: 0.4, B: 0.2, A: 1},
			{R: 0.95, G: 0.7, B: 0.2, A: 1},
			{R: 0.7, G: 0.2, B: 0.2, A: 1},
		}
	case 2: // cool
		return []core.RGBA{
			{R: 0.2, G: 0.5, B: 0.9, A: 1},
			{R: 0.3, G: 0.8, B: 0.7, A: 1},
			{R: 0.1, G: 0.2, B: 0.6, A: 1},
		}
	default: // mono
		return []core.RGBA{{R: 0, G: 0, B: 0, A: 1}}
	}
}
