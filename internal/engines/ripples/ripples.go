// Package ripples draws the interference of circular waves from a few
// sources: rings that cross into moiré lattices.
//
// All parameters are continuous. Sources sit on fixed circles chosen by the
// seed (always the same random draws), drift moves them along those circles
// and phase moves the waves outward; phase loops at 1, drift at 1.
package ripples

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
)

type Engine struct{}

func (Engine) Name() string { return "ripples" }

const maxSources = 4

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	sources := int(core.Pick(params, "sources", 3)) // 1..4 (a count: keep it fixed when sweeping)
	freq := core.Pick(params, "freq", 9)            // rings per unit distance
	phase := core.Pick(params, "phase", 0)          // waves move out one ring per unit, loops
	drift := core.Pick(params, "drift", 0)          // sources move round their orbits, loops at 1
	orbit := core.Pick(params, "orbit", 0.25)       // orbit radius
	sharp := core.Pick(params, "sharp", 3)          // 0 = soft gradients, higher = crisp bands
	falloff := core.Pick(params, "falloff", 0.6)    // how fast a source's waves fade with distance

	// Always draw maxSources orbits so the sources count doesn't shift the
	// rest of the random stream.
	type src struct{ cx, cy, a0, dir float64 }
	var ss [maxSources]src
	for i := range ss {
		ss[i] = src{0.5 + 0.15*(rng.Float64()-0.5), 0.5 + 0.15*(rng.Float64()-0.5), 2 * math.Pi * rng.Float64(), 1}
		if rng.Float64() < 0.5 {
			ss[i].dir = -1
		}
	}
	sources = min(max(sources, 1), maxSources)

	n := len(colors)
	scene := core.Scene{}
	scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
		sum, wsum := 0.0, 0.0
		for i := 0; i < sources; i++ {
			s := ss[i]
			a := s.a0 + s.dir*2*math.Pi*drift
			px, py := s.cx+orbit*math.Cos(a), s.cy+orbit*math.Sin(a)
			d := math.Hypot(x-px, y-py)
			w := math.Exp(-falloff * d)
			sum += w * math.Cos(2*math.Pi*(d*freq-phase))
			wsum += w
		}
		v := sum / wsum // -1..1
		if sharp > 0 {
			v = math.Tanh(sharp*v) / math.Tanh(sharp)
		}
		t := 0.5 + 0.5*v
		return colors[min(n-1, int(t*float64(n)))]
	}})
	return scene, nil
}
