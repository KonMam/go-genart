package swirl

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/colorize"
	"genart/internal/core"
	"genart/internal/geom"
	"genart/internal/noise"
	"genart/internal/randutil"
)

type Engine struct{}

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

func (Engine) Name() string { return "swirl" }

func (Engine) Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// Parameters
	circleN := int(core.Pick(params, "circles", 500))
	dotsN := int(core.Pick(params, "dots", 100))
	lineWidth := core.Pick(params, "lw", 0.001)
	nIters := int(core.Pick(params, "nIters", 1000))
	factor := core.Pick(params, "factor", 1.5)
	step := core.Pick(params, "step", 0.003)
	maxRadius := core.Pick(params, "maxRadius", 0.05)

	scene := core.Scene{}

	// generate circles in a spiral
	cs := make([]circles, 0)
	goldenAngle := math.Pi * (3.0 - math.Sqrt(5.0))
	for i := 0; i < circleN; i++ {
		theta := goldenAngle * float64(i)
		r := math.Sqrt(float64(i)/float64(circleN)) * 0.45
		cs = append(cs, circles{
			x:      0.5 + r*math.Cos(theta),
			y:      0.5 + r*math.Sin(theta),
			radius: randutil.RandomRangeFloat64(rng, maxRadius*0.1, maxRadius),
		})
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
	const epsilon = 0.001

	for i := 0; i < circleN; i++ {
		for j := 0; j < nIters; j++ {
			for k := range ds[i] {
				// curl noise
				nx, ny := noise.Curl2D(noiseField, ds[i][k].x*factor, ds[i][k].y*factor, epsilon)

				ds[i][k].prevx, ds[i][k].prevy = ds[i][k].x, ds[i][k].y
				ds[i][k].x += nx * ds[i][k].step
				ds[i][k].y += ny * ds[i][k].step

				// inside the stroke drawing loop
				// pick color based on noise value at current position
				c := colorize.PickColorFromNoise(colors, noiseField, ds[i][k].x, ds[i][k].y, factor)

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

		
	}

	return scene, nil
}


