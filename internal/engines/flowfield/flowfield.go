package flowfield

import (
	"context"
	"fmt"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/noise"
)

type Engine struct{}

func (Engine) Name() string { return "flowfield" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// --- Params with defaults ---
	particles := int(pick(params, "particles", 1000))
	steps := int(pick(params, "steps", 300))
	scale := pick(params, "scale", 0.002)
	step := pick(params, "step", 0.002)
	lw := pick(params, "lw", 0.0015)

	if particles <= 0 {
		return core.Scene{}, fmt.Errorf("invalid particles %d (must be > 0)", particles)
	}
	if steps <= 0 {
		return core.Scene{}, fmt.Errorf("invalid steps %d (must be > 0)", steps)
	}
	if scale <= 0 {
		return core.Scene{}, fmt.Errorf("invalid scale %f (must be > 0)", scale)
	}
	if step <= 0 {
		return core.Scene{}, fmt.Errorf("invalid step %f (must be > 0)", step)
	}
	if lw < 0 {
		return core.Scene{}, fmt.Errorf("invalid lw %f (must be >= 0)", lw)
	}

	// --- Field: deterministic with sub-seed ---
	field := noise.NewSimplexField(rng.Int63(), scale)

	scene := core.Scene{}

	for i := 0; i < particles; i++ {
		// random start in [0,1]
		x, y := rng.Float64(), rng.Float64()
		points := make([]core.Vec2, 0, steps)

		for j := 0; j < steps; j++ {
			// field → angle in radians
			angle := field.At(x, y) * math.Pi // [-1,1] → [-π,π]
			dx := math.Cos(angle) * step
			dy := math.Sin(angle) * step

			x += dx
			y += dy

			// stop if out of bounds
			if x < 0 || x > 1 || y < 0 || y > 1 {
				break
			}

			points = append(points, core.Vec2{X: x, Y: y})
		}

		if len(points) > 1 {
			scene.Items = append(scene.Items, core.Stroke{
				Path:  core.Path{Points: points, Closed: false},
				Width: lw,
				Color: core.RGBA{R: 0, G: 0, B: 0, A: 0.3}, // translucent black
				Alpha: 0.3,
			})
		}
	}

	return scene, nil
}

func pick(m map[string]float64, k string, def float64) float64 {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}
