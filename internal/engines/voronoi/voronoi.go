// Package voronoi draws stained glass: Voronoi cells around seeds that
// drift on closed orbits, each cell shaded like a pebble (light at its
// seed, darker towards its edges) with dark lead lines between cells.
//
// All parameters but cols and rows are continuous; phase loops at 1, since every
// seed's orbit is a Lissajous figure with whole-number frequencies. The
// engine draws the same random numbers whatever the parameters.
package voronoi

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
)

type Engine struct{}

func (Engine) Name() string { return "voronoi" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	cols := int(core.Pick(params, "cols", 6)) // seeds sit on a jittered cols x rows grid (counts)
	rows := int(core.Pick(params, "rows", 0)) // 0: as many as keep cells about square

	phase := core.Pick(params, "phase", 0)      // position on the orbits, loops at 1
	wander := core.Pick(params, "wander", 0.06) // orbit size
	lead := core.Pick(params, "lead", 0.012)    // width of the lines between cells
	dome := core.Pick(params, "dome", 1)        // 0 = flat cells, 1 = shaded like pebbles
	aspect := core.Pick(params, "aspect", 160.0/144)

	// One seed per cell of a jittered grid (with a ring of cells outside
	// the frame so edge cells close), so seeds never pile up into slivers.
	// Each draws the same six numbers whatever the continuous parameters.
	cols = max(cols, 1)
	if rows <= 0 {
		rows = int(math.Round(float64(cols) / aspect))
	}
	rows = max(rows, 1)
	cw, ch := aspect/float64(cols), 1/float64(rows)
	a := 2 * math.Pi * phase
	type pt struct{ x, y, tone float64 }
	var ps []pt
	for j := -1; j <= rows; j++ {
		for i := -1; i <= cols; i++ {
			jx, jy := rng.Float64(), rng.Float64()
			fx, fy := float64(1+rng.Intn(2)), float64(1+rng.Intn(2))
			px, py := 2*math.Pi*rng.Float64(), 2*math.Pi*rng.Float64()
			tone := rng.Float64()
			x := (float64(i) + 0.2 + 0.6*jx) * cw
			y := (float64(j) + 0.2 + 0.6*jy) * ch
			ps = append(ps, pt{x + wander*math.Sin(fx*a+px), y + wander*math.Sin(fy*a+py), tone})
		}
	}

	n := len(colors)
	scene := core.Scene{}
	scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
		x *= aspect
		d1, d2 := math.Inf(1), math.Inf(1)
		best, second := 0, 0
		for i, p := range ps {
			d := (x-p.x)*(x-p.x) + (y-p.y)*(y-p.y)
			if d < d1 {
				d2, second = d1, best
				d1, best = d, i
			} else if d < d2 {
				d2, second = d, i
			}
		}
		// Exact distance to the bisector of the two nearest seeds.
		p, q := ps[best], ps[second]
		edge := (d2 - d1) / (2 * math.Hypot(q.x-p.x, q.y-p.y))
		d1 = math.Sqrt(d1)
		if edge < lead/2 {
			return colors[0]
		}
		// Brighter towards the middle of the cell; base tone by seed.
		shape := math.Min(1, edge/(d1+edge)*2)
		t := 0.3 + 0.7*((1-dome)+dome*shape)*(0.55+0.45*ps[best].tone)
		return colors[min(n-1, 1+int(t*float64(n-1)))]
	}})
	return scene, nil
}
