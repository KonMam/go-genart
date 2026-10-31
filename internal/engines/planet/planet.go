// Package planet draws a lit, spinning planet: continents and seas from 3D
// noise on the sphere, a day side and a night side, and a thin ring of
// atmosphere.
//
// All parameters are continuous. spin turns the planet about its axis and
// loops at 1 (the noise is sampled on the rotated sphere, so the surface is
// the same after a full turn); light moves the sun round it. The engine
// draws one random number whatever the parameters.
package planet

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"github.com/ojrac/opensimplex-go"
)

type Engine struct{}

func (Engine) Name() string { return "planet" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	spin := core.Pick(params, "spin", 0)             // rotation, loops at 1
	light := core.Pick(params, "light", 0.15)        // sun direction round the planet, loops at 1
	tilt := core.Pick(params, "tilt", 0.35)          // axis tilt, radians
	sea := core.Pick(params, "sea", 0.0)             // sea level, -1..1
	freq := core.Pick(params, "freq", 1.8)           // continents per radius
	radius := core.Pick(params, "radius", 0.4)       // planet size
	aspect := core.Pick(params, "aspect", 160.0/144) // output width/height, so the planet stays round

	nz := opensimplex.New(rng.Int63())
	n := len(colors)
	a := 2 * math.Pi * spin
	la := 2 * math.Pi * light
	lx, ly, lz := math.Sin(la), -0.35, math.Cos(la)
	ll := math.Sqrt(lx*lx + ly*ly + lz*lz)
	lx, ly, lz = lx/ll, ly/ll, lz/ll

	scene := core.Scene{}
	scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
		dx, dy := (x-0.5)*aspect/radius, (y-0.5)/radius
		r2 := dx*dx + dy*dy
		if r2 > 1 {
			// Atmosphere: a thin lit rim just outside the disc.
			if r := math.Sqrt(r2); r < 1.06 {
				lit := (dx*lx + dy*ly) / r
				if lit > 0.2 {
					return colors[1]
				}
			}
			return colors[0]
		}
		dz := math.Sqrt(1 - r2) // towards the viewer
		// Light the visible point.
		shade := dx*lx + dy*ly + dz*lz
		// Into planet coordinates: undo the tilt, then the spin.
		py := dy*math.Cos(tilt) - dz*math.Sin(tilt)
		pz := dy*math.Sin(tilt) + dz*math.Cos(tilt)
		px := dx
		qx := px*math.Cos(a) + pz*math.Sin(a)
		qz := -px*math.Sin(a) + pz*math.Cos(a)
		h := nz.Eval3(qx*freq, py*freq, qz*freq) + 0.4*nz.Eval3(qx*freq*2.3+9, py*freq*2.3, qz*freq*2.3)
		land := h > sea
		// Brightness: sea darker than land, both scaled by the light.
		base := 0.45
		if land {
			base = 1
		}
		// Light in three steps, like a cel-shaded globe.
		t := base * math.Max(0, 0.15+0.85*math.Round(math.Max(0, shade)*3)/3)
		return colors[min(n-1, 1+int(t*float64(n-1)))]
	}})
	return scene, nil
}
