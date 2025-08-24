package jittersquare

import (
	"context"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
)

type Engine struct{}

func (Engine) Name() string { return "jittersquare" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64) (core.Scene, error) {
	// defaults
	size := pick(params, "size", 0.6)
	lw := pick(params, "lw", 0.003)
	jitter := pick(params, "jitter", 0.02)

	// base square
	path := geom.Square(0.5, 0.5, size)

	// jitter vertices
	jittered := make([]core.Vec2, len(path.Points))
	for i, p := range path.Points {
		jx := (rng.Float64()*2 - 1) * jitter // uniform [-jitter, +jitter]
		jy := (rng.Float64()*2 - 1) * jitter
		jittered[i] = core.Vec2{X: p.X + jx, Y: p.Y + jy}
	}
	jpath := core.Path{Points: jittered, Closed: true}

	scene := core.Scene{
		Items: []core.Item{
			core.Fill{
				Polygon: jpath,
				Color:   core.RGBA{R: 0.8, G: 0.4, B: 0.2, A: 0.9}, // orange-ish
				Alpha:   0.9,
			},
			core.Stroke{
				Path:  jpath,
				Width: lw,
				Color: core.RGBA{R: 0, G: 0, B: 0, A: 1},
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
