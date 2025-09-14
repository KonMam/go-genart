package flow

import (
	"context"
	"math/rand"

	"genart/internal/colorize"
	"genart/internal/core"
	"genart/internal/noise"
	"genart/internal/randutil"
)

type Engine struct{}

type dot struct {
	x, y         float64
	prevx, prevy float64
}

func (Engine) Name() string { return "flow" }

func (Engine) Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// Parameters
	dotsN := int(core.Pick(params, "dots", 5000))
	lineWidth := core.Pick(params, "lw", 0.001)
	nIters := int(core.Pick(params, "nIters", 100))
	factor := core.Pick(params, "factor", 1.5)
	step := core.Pick(params, "step", 0.005)

	scene := core.Scene{}

	// initialize dots
	ds := make([]dot, 0)
	for i := 0; i < dotsN; i++ {
		x := rng.Float64()
		y := rng.Float64()
		ds = append(ds, dot{
			x:     x,
			y:     y,
			prevx: x,
			prevy: y,
		})
	}

	noiseField := noise.NewPerlinField(rng.Int63(), 1.0)
	const epsilon = 0.001

	for i := 0; i < nIters; i++ {
		for k := range ds {
			// curl noise
			nx, ny := noise.Curl2D(noiseField, ds[k].x*factor, ds[k].y*factor, epsilon)

			ds[k].prevx, ds[k].prevy = ds[k].x, ds[k].y
			ds[k].x += nx * step
			ds[k].y += ny * step

			// pick color based on noise value at current position
			c := colorize.PickColorFromNoise(colors, noiseField, ds[k].x, ds[k].y, factor)

			alpha := 0.05 + rng.Float64()*0.1
			lw := lineWidth * randutil.RandomRangeFloat64(rng, 0.8, 1.2)
			points := []core.Vec2{
				{X: ds[k].prevx, Y: ds[k].prevy},
				{X: ds[k].x, Y: ds[k].y},
			}
			scene.AddStroke(points, false, lw, c, alpha)

			// wrap around
			if ds[k].x < 0 {
				ds[k].x = 1
				ds[k].prevx = ds[k].x
			}
			if ds[k].x > 1 {
				ds[k].x = 0
				ds[k].prevx = ds[k].x
			}
			if ds[k].y < 0 {
				ds[k].y = 1
				ds[k].prevy = ds[k].y
			}
			if ds[k].y > 1 {
				ds[k].y = 0
				ds[k].prevy = ds[k].y
			}
		}
	}

	return scene, nil
}


