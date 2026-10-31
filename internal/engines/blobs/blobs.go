// Package blobs draws metaballs: soft blobs that drift on closed orbits and
// melt into each other, with a glowing halo shaded in bands.
//
// All parameters are continuous. Each blob's orbit is a Lissajous figure
// with whole-number frequencies chosen by the seed, so phase 1 brings every
// blob back to where it started (phase loops). The engine draws the same
// random numbers whatever the parameters.
package blobs

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
)

type Engine struct{}

func (Engine) Name() string { return "blobs" }

const maxBlobs = 12

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	count := int(core.Pick(params, "blobs", 7)) // a count: keep it fixed when sweeping
	radius := core.Pick(params, "radius", 0.09) // blob size
	spread := core.Pick(params, "spread", 0.3)  // orbit size
	phase := core.Pick(params, "phase", 0)      // position on the orbits, loops at 1
	glow := core.Pick(params, "glow", 0.5)      // halo width: 0 = hard edges
	rings := core.Pick(params, "rings", 0)      // contour rings in the halo (0 = smooth gradient)

	type blob struct{ fx, fy, px, py, r float64 }
	var bs [maxBlobs]blob
	for i := range bs {
		bs[i] = blob{
			fx: float64(1 + rng.Intn(3)), fy: float64(1 + rng.Intn(3)),
			px: 2 * math.Pi * rng.Float64(), py: 2 * math.Pi * rng.Float64(),
			r: 0.6 + 0.8*rng.Float64(),
		}
	}
	count = min(max(count, 1), maxBlobs)

	n := len(colors)
	scene := core.Scene{}
	scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
		f := 0.0
		for i := 0; i < count; i++ {
			b := bs[i]
			a := 2 * math.Pi * phase
			cx := 0.5 + spread*math.Sin(b.fx*a+b.px)
			cy := 0.5 + spread*math.Sin(b.fy*a+b.py)
			r := radius * b.r
			d2 := (x-cx)*(x-cx) + (y-cy)*(y-cy)
			f += r * r / (d2 + 1e-9)
		}
		// f = 1 is the blob surface. Inside is solid; outside the halo
		// fades over glow.
		var t float64
		if f >= 1 {
			t = 1
		} else if glow > 0 {
			t = math.Max(0, 1-(1-f)/glow) * 0.75
			if rings > 0 {
				t = math.Floor(t*rings) / rings
			}
		}
		return colors[min(n-1, int(t*float64(n)))]
	}})
	return scene, nil
}
