package perlinpearls

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
	"genart/internal/noise"
	"genart/internal/randutil"
)

type Engine struct{}

func (Engine) Name() string { return "perlinpearls" }

func (Engine) Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// Parameters
	circleN := int(pick(params, "circles", 5))
	dotsN := int(pick(params, "dots", 500))
	lineWidth := pick(params, "lw", 0.001)
	nIters := int(pick(params, "nIters", 2000))
	factor := pick(params, "factor", 1.5)
	step := pick(params, "step", 0.003)

	// outline params
	outlineWidth := pick(params, "outlineWidth", lineWidth*2)
	scene := core.Scene{}

	// generate non-overlapping circles
	cs := make([]circles, 0)
	for len(cs) < circleN {
		c := circles{
			x:      randutil.RandomRangeFloat64(rng, 0.1, 0.9),
			y:      randutil.RandomRangeFloat64(rng, 0.1, 0.9),
			radius: randutil.RandomRangeFloat64(rng, 0.05, 0.2),
		}
		var overlapping bool
		for _, cl := range cs {
			d := (geom.Vec2{X: c.x, Y: c.y}).Distance(geom.Vec2{X: cl.x, Y: cl.y})
			if d < c.radius+cl.radius {
				overlapping = true
				break
			}
		}
		if !overlapping {
			cs = append(cs, c)
		}
	}

	// initialize dots around each circle
	ds := make([][]dot, 0)
	for i := 0; i < circleN; i++ {
		dots := make([]dot, 0)
		for j := 0; j < dotsN; j++ {
			theta := rng.Float64() * math.Pi * 2
			dotStep := step * randutil.RandomRangeFloat64(rng, 0.8, 1.2)
			dots = append(dots, dot{
				theta: theta,
				cx:    cs[i].x,
				cy:    cs[i].y,
				x:     cs[i].x + math.Sin(theta)*cs[i].radius,
				y:     cs[i].y + math.Cos(theta)*cs[i].radius,
				prevx: cs[i].x + math.Sin(theta)*cs[i].radius,
				prevy: cs[i].y + math.Cos(theta)*cs[i].radius,
				step:  dotStep,
			})
		}
		ds = append(ds, dots)
	}

	noiseField := noise.NewPerlinField(rng.Int63(), 1.0)
	eps := 0.001

	for i := 0; i < circleN; i++ {
		for j := 0; j < nIters; j++ {
			for k := range ds[i] {
				// curl noise
				n1 := noiseField.At((ds[i][k].x+eps)*factor, ds[i][k].y*factor)
				n2 := noiseField.At((ds[i][k].x-eps)*factor, ds[i][k].y*factor)
				n3 := noiseField.At(ds[i][k].x*factor, (ds[i][k].y+eps)*factor)
				n4 := noiseField.At(ds[i][k].x*factor, (ds[i][k].y-eps)*factor)

				dx := (n1 - n2) / (2 * eps)
				dy := (n3 - n4) / (2 * eps)

				nx := dy
				ny := -dx

				ds[i][k].prevx, ds[i][k].prevy = ds[i][k].x, ds[i][k].y
				ds[i][k].x += nx * ds[i][k].step
				ds[i][k].y += ny * ds[i][k].step

				// inside the stroke drawing loop
				// pick color based on noise value at current position
				var c core.RGBA
				if len(colors) > 0 {
					nval := noiseField.At(ds[i][k].x*factor, ds[i][k].y*factor) // -1..1
					t := (nval + 1) / 2                                         // normalize 0..1
					idx := int(t * float64(len(colors)-1))
					if idx < 0 {
						idx = 0
					}
					if idx >= len(colors) {
						idx = len(colors) - 1
					}
					c = colors[idx]
				} else {
					c = core.RGBA{0, 0, 0, 1}
				}

				// only draw if inside circle
				if (geom.Vec2{X: ds[i][k].x, Y: ds[i][k].y}).Distance(geom.Vec2{X: cs[i].x, Y: cs[i].y}) < cs[i].radius &&
					(geom.Vec2{X: ds[i][k].prevx, Y: ds[i][k].prevy}).Distance(geom.Vec2{X: cs[i].x, Y: cs[i].y}) < cs[i].radius {
					alpha := 0.05 + rng.Float64()*0.1
					lw := lineWidth * randutil.RandomRangeFloat64(rng, 0.8, 1.2)
					points := []core.Vec2{
						{X: ds[i][k].prevx, Y: ds[i][k].prevy},
						{X: ds[i][k].x, Y: ds[i][k].y},
					}
					scene.AddStroke(points, false, lw, c, alpha)
				}
			}
		}

		// --- Circle outline (always black now) ---
		segments := 200
		outlinePoints := make([]core.Vec2, 0, segments+1)
		for s := 0; s <= segments; s++ {
			theta := 2 * math.Pi * float64(s) / float64(segments)
			x := cs[i].x + math.Cos(theta)*cs[i].radius
			y := cs[i].y + math.Sin(theta)*cs[i].radius
			outlinePoints = append(outlinePoints, core.Vec2{X: x, Y: y})
		}

		outlineColor := core.RGBA{R: 0, G: 0, B: 0, A: 1} // black outline
		scene.AddStroke(outlinePoints, true, outlineWidth, outlineColor, 1.0)
	}

	return scene, nil
}

type circles struct {
	x, y   float64
	radius float64
}

type dot struct {
	theta        float64
	cx, cy       float64
	x, y         float64
	prevx, prevy float64
	step         float64
}

func pick(m map[string]float64, k string, def float64) float64 {
	if v, ok := m[k]; ok {
		return v
	}
	return def
}
