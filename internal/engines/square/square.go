package square

import (
	"context"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
)

type Engine struct{}

func (Engine) Name() string { return "square" }

func (Engine) Generate(_ context.Context, _ *rand.Rand, params map[string]float64) (core.Scene, error) {
	// defaults
	size := pick(params, "size", 0.6)
	lw := pick(params, "lw", 0.003)

	path := geom.Square(0.5, 0.5, size)

	scene := core.Scene{
		Items: []core.Item{
			core.Fill{
				Polygon: path,
				Color:   core.RGBA{R: 0.9, G: 0.2, B: 0.2, A: 1}, // reddish fill
				Alpha:   1,
			},
			core.Stroke{
				Path:  path,
				Width: lw,
				Color: core.RGBA{R: 0.1, G: 0.1, B: 0.1, A: 1}, // dark stroke
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
