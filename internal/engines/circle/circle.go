package circle

import (
	"context"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
)

type Engine struct{}

func (Engine) Name() string { return "circle" }

func (Engine) Generate(_ context.Context, _ *rand.Rand, params map[string]float64) (core.Scene, error) {
	// defaults
	radius := pick(params, "radius", 0.3)
	segments := int(pick(params, "segments", 128))
	if segments < 8 {
		segments = 8
	}
	lw := pick(params, "lw", 0.003)

	path := geom.Circle(0.5, 0.5, radius, segments)

	scene := core.Scene{
		Items: []core.Item{
			core.Fill{
				Polygon: path,
				Color:   core.RGBA{R: 0.2, G: 0.5, B: 0.9, A: 0.9}, // bluish fill
				Alpha:   0.9,
			},
			core.Stroke{
				Path:  path,
				Width: lw,
				Color: core.RGBA{R: 0.05, G: 0.05, B: 0.05, A: 1}, // dark stroke
				Alpha: 1,
			},
		},
	}
	return scene, nil
}

func pick(m map[string]float64, k string, def float64) float64 {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}
