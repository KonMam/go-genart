package blackhole

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/noise"
	"genart/internal/palette"
)

type Engine struct{}

func (Engine) Name() string { return "blackhole" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64) (core.Scene, error) {
	// Parameters
	circleN := int(pick(params, "circles", 120))
	density := pick(params, "density", 0.6)
	circleGap := pick(params, "gap", 0.02)
	lineWidth := pick(params, "lw", 0.0008)
	segments := int(pick(params, "segments", 900))
	hole := pick(params, "hole", 0.1)
	paletteID := int(pick(params, "palette", 1))
	freq := pick(params, "freq", 6.0)
	amp := pick(params, "amp", 1.2)

	centerX, centerY := 0.5, 0.5
	radiusOuter := 0.45

	// Noise field
	field := noise.NewSimplexField3D(rng.Int63(), 1.0)
	colors := selectPalette(paletteID)

	scene := core.Scene{}

	kMax := 0.5 + rng.Float64()*0.5

	for i := 0; i < circleN; i++ {
		t := float64(i) / float64(circleN)
		radius := hole + t*(radiusOuter-hole)

		k := kMax * math.Sqrt(t)
		noisiness := density * t * t

		points := make([]core.Vec2, 0, segments)

		// NEW: random starting offset angle
		startTheta := rng.Float64() * 2 * math.Pi

		for j := 0; j < segments; j++ {
			theta := startTheta + 2*math.Pi*float64(j)/float64(segments)

			// High-frequency noise
			r1 := math.Cos(theta) + 1
			r2 := math.Sin(theta) + 1
			nv := field.At(k*freq*r1, k*freq*r2, float64(i)*circleGap)

			r := radius + nv*noisiness*amp
			if r < hole {
				r = hole
			}

			x := centerX + r*math.Cos(theta)
			y := centerY + r*math.Sin(theta)

			points = append(points, core.Vec2{X: x, Y: y})
		}

		// NEW: alpha jitter to reduce banding
		alpha := 0.6 + rng.Float64()*0.25
		scene.AddStroke(points, true, lineWidth, colors.Pick(rng), alpha)
	}

	return scene, nil
}

// --- helpers ---

func pick(m map[string]float64, k string, def float64) float64 {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}

func selectPalette(id int) palette.Palette {
	switch id {
	case 1:
		return palette.Warm
	case 2:
		return palette.Cool
	case 3:
		return palette.Rainbow
	default:
		return palette.Mono
	}
}
