package strata

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
	"genart/internal/noise"
)

type Engine struct{}

func (Engine) Name() string { return "strata" }

func (Engine) Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// Parameters
	sides := int(core.Pick(params, "sides", 6))
	layers := int(core.Pick(params, "layers", 20))
	depth := int(core.Pick(params, "depth", 5))
	magnitude := core.Pick(params, "magnitude", 0.1)
	rotation := core.Pick(params, "rotation", 0.01)

	scene := core.Scene{}
	noiseField := noise.NewPerlinField(rng.Int63(), 1.0)

	for i := 0; i < layers; i++ {
		// Base polygon
		points := make([]geom.Vec2, 0, sides)
		radius := 0.5 - float64(i)*0.02
		if radius < 0 {
			radius = 0
		}
		angleOffset := float64(i) * rotation
		for j := 0; j < sides; j++ {
			angle := float64(j)/float64(sides)*2*math.Pi + angleOffset
			points = append(points, geom.Vec2{
				X: 0.5 + radius*math.Cos(angle),
				Y: 0.5 + radius*math.Sin(angle),
			})
		}

		// Subdivide and displace
		finalGeomPoints := subdivide(points, depth, magnitude, noiseField)

		// Convert to core.Vec2
		finalCorePoints := make([]core.Vec2, len(finalGeomPoints))
		for i, p := range finalGeomPoints {
			finalCorePoints[i] = core.Vec2{X: p.X, Y: p.Y}
		}

		// Pick color
		color := colors[i%len(colors)]

		// Add to scene
		scene.AddFill(finalCorePoints, color, 0.8)
	}

	return scene, nil
}

func subdivide(points []geom.Vec2, depth int, magnitude float64, noiseField noise.ScalarField2D) []geom.Vec2 {
	if depth == 0 {
		return points
	}

	newPoints := make([]geom.Vec2, 0)

	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i+1)%len(points)]

		mid := geom.Vec2{
			X: (p1.X + p2.X) / 2,
			Y: (p1.Y + p2.Y) / 2,
		}

		dx := p2.X - p1.X
		dy := p2.Y - p1.Y

		normal := geom.Vec2{X: -dy, Y: dx}
		normal = normal.Normalize()

		noiseVal := noiseField.At(mid.X, mid.Y)
		displacement := normal.Scale(noiseVal * magnitude)

		mid = mid.Add(displacement)

		newPoints = append(newPoints, p1, mid)
	}

	return subdivide(newPoints, depth-1, magnitude, noiseField)
}
