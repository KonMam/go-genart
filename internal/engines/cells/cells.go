package cells

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

func (Engine) Name() string { return "cells" }

func (Engine) Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	// Parameters
	gridSize := int(core.Pick(params, "grid_size", 10))
	dotsPerCell := int(core.Pick(params, "dots_per_cell", 500))
	lineWidth := core.Pick(params, "lw", 0.001)
	nIters := int(core.Pick(params, "nIters", 100))
	factor := core.Pick(params, "factor", 1.5)
	step := core.Pick(params, "step", 0.005)

	scene := core.Scene{}

	// initialize dots in each cell
	cells := make([][]*dot, gridSize*gridSize)
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			cellIndex := i*gridSize + j
			cells[cellIndex] = make([]*dot, dotsPerCell)

			cellWidth := 1.0 / float64(gridSize)
			cellHeight := 1.0 / float64(gridSize)
			cellX := float64(j) * cellWidth
			cellY := float64(i) * cellHeight

			for k := 0; k < dotsPerCell; k++ {
				x := cellX + rng.Float64()*cellWidth
				y := cellY + rng.Float64()*cellHeight
				cells[cellIndex][k] = &dot{
					x:     x,
					y:     y,
					prevx: x,
					prevy: y,
				}
			}
		}
	}

	noiseField := noise.NewPerlinField(rng.Int63(), 1.0)
	const epsilon = 0.001

	cellWidth := 1.0 / float64(gridSize)
	cellHeight := 1.0 / float64(gridSize)

	for i := 0; i < nIters; i++ {
		for j := 0; j < len(cells); j++ {
			cellX := float64(j%gridSize) * cellWidth
			cellY := float64(j/gridSize) * cellHeight

			for k := range cells[j] {
				dot := cells[j][k]

				// curl noise
				nx, ny := noise.Curl2D(noiseField, dot.x*factor, dot.y*factor, epsilon)

				dot.prevx, dot.prevy = dot.x, dot.y
				dot.x += nx * step
				dot.y += ny * step

				// pick color based on noise value at current position
				c := colorize.PickColorFromNoise(colors, noiseField, dot.x, dot.y, factor)

				alpha := 0.05 + rng.Float64()*0.1
				lw := lineWidth * randutil.RandomRangeFloat64(rng, 0.8, 1.2)
				points := []core.Vec2{
					{X: dot.prevx, Y: dot.prevy},
					{X: dot.x, Y: dot.y},
				}
				scene.AddStroke(points, false, lw, c, alpha)

				// wrap around cell
				if dot.x < cellX {
					dot.x = cellX + cellWidth
					dot.prevx = dot.x
				}
				if dot.x > cellX+cellWidth {
					dot.x = cellX
					dot.prevx = dot.x
				}
				if dot.y < cellY {
					dot.y = cellY + cellHeight
					dot.prevy = dot.y
				}
				if dot.y > cellY+cellHeight {
					dot.y = cellY
					dot.prevy = dot.y
				}
			}
		}
	}

	return scene, nil
}
