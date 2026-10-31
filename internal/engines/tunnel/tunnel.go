// Package tunnel draws the demoscene tunnel: a checkered tube seen from
// inside, receding to a dark vanishing point that wanders.
//
// All parameters are continuous. depth moves the camera forward and loops
// at 1 (one checker length), turn spins the tube and loops at 1. The engine
// draws no random numbers.
package tunnel

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
)

type Engine struct{}

func (Engine) Name() string { return "tunnel" }

func (Engine) Generate(_ context.Context, _ *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	depth := core.Pick(params, "depth", 0)      // camera position along the tube, loops at 1
	turn := core.Pick(params, "turn", 0)        // tube rotation, loops at 1
	sectors := core.Pick(params, "sectors", 12) // checker columns round the tube (keep it whole to avoid a seam)
	sway := core.Pick(params, "sway", 0.12)     // how far the vanishing point leans off centre
	bend := core.Pick(params, "bend", 0)        // direction of the lean, loops at 1
	fog := core.Pick(params, "fog", 0.45)       // radius at which the far end fades to black
	aspect := core.Pick(params, "aspect", 160.0/144)

	n := len(colors)
	ba := 2 * math.Pi * bend
	vx, vy := 0.5+sway*math.Cos(ba), 0.5+sway*math.Sin(ba)
	scene := core.Scene{}
	scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
		// Lean the centre towards the vanishing point for far (small
		// radius) points, so the tube curves away.
		dx, dy := (x-0.5)*aspect, y-0.5
		r := math.Hypot(dx, dy)
		k := math.Exp(-r / 0.15)
		dx, dy = (x-(0.5+(vx-0.5)*k))*aspect, y-(0.5+(vy-0.5)*k)
		r = math.Max(math.Hypot(dx, dy), 1e-6)
		z := 0.6/r + depth*2 // two rings per checker
		a := math.Atan2(dy, dx)/(2*math.Pi) + turn
		check := (int(math.Floor(z))+int(math.Floor(a*sectors)))&1 == 0
		light := math.Min(1, r/fog) // far = dark
		t := light
		if !check {
			t *= 0.55
		}
		return colors[min(n-1, int(t*float64(n)))]
	}})
	return scene, nil
}
