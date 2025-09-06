package randutil

import (
	"math/rand"

	"genart/internal/geom"
)

// RandomRangeFloat64 returns a random float64 in [min, max).
func RandomRangeFloat64(rng *rand.Rand, min, max float64) float64 {
	return min + rng.Float64()*(max-min)
}

// RandomRangeFloat64 returns a random float64 in [min, max).
func RandomRangeInt(rng *rand.Rand, min, max int) int {
	return min + rng.Int()*(max-min)
}

// UniformPoint returns a random point in [0,1] × [0,1].
func UniformPoint(rng *rand.Rand) geom.Vec2 {
	return geom.Vec2{X: rng.Float64(), Y: rng.Float64()}
}

// JitteredGrid returns n*n points arranged on a jittered grid inside [0,1] × [0,1].
//
// gridSize controls subdivision (e.g. 10 → 100 cells).
// jitter is in [0,1], 0 = perfect grid, 1 = max jitter (point anywhere in cell).
func JitteredGrid(rng *rand.Rand, gridSize int, jitter float64) []geom.Vec2 {
	points := make([]geom.Vec2, 0, gridSize*gridSize)
	cellSize := 1.0 / float64(gridSize)

	for gx := 0; gx < gridSize; gx++ {
		for gy := 0; gy < gridSize; gy++ {
			// center of the cell
			cx := (float64(gx) + 0.5) * cellSize
			cy := (float64(gy) + 0.5) * cellSize

			// jitter offset
			jx := (rng.Float64() - 0.5) * jitter * cellSize
			jy := (rng.Float64() - 0.5) * jitter * cellSize

			points = append(points, geom.Vec2{X: cx + jx, Y: cy + jy})
		}
	}

	return points
}
