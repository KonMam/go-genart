// Package truchet draws a grid of Truchet tiles, two quarter-circle arcs
// per tile, whose arcs join into long meandering paths. Each tile turns
// smoothly between its two orientations as a periodic noise field
// (noise.LoopField) passes over it, so the paths unravel and re-knit.
//
// All parameters but cells are continuous; phase loops. The engine draws one
// random number whatever the parameters.
package truchet

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/geom"
	"genart/internal/noise"
)

type Engine struct{}

func (Engine) Name() string { return "truchet" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	cells := int(core.Pick(params, "cells", 8)) // tiles across, a count
	phase := core.Pick(params, "phase", 0)      // loops at 1
	travel := core.Pick(params, "travel", 0.6)  // how much one phase cycle turns the tiles
	scale := core.Pick(params, "scale", 0.35)   // size of regions that turn together
	snap := core.Pick(params, "snap", 4)        // 0 = tiles drift at any angle, higher = they rest at 0 or 90 degrees
	weight := core.Pick(params, "weight", 0.22) // arc width as a fraction of a tile

	field := noise.NewLoopField(rng.Int63(), scale, travel)
	fg := colors[len(colors)-1]
	dim := colors[len(colors)/2]
	s := 1 / float64(cells)
	lw := weight * s
	scene := core.Scene{}
	rowsN := int(math.Ceil(float64(cells) * 1.0))
	for j := 0; j < rowsN; j++ {
		for i := 0; i < cells; i++ {
			cx, cy := (float64(i)+0.5)*s, (float64(j)+0.5)*s
			v := 0.5 + 0.5*math.Tanh(snap*field.At(cx, cy, phase)) // 0..1
			ang := v * math.Pi / 2
			// Tiles caught mid-turn draw dimmer.
			c := fg
			if math.Abs(v-0.5) < 0.3 {
				c = dim
			}
			for _, corner := range [2][2]float64{{-1, -1}, {1, 1}} {
				ox, oy := corner[0]*s/2, corner[1]*s/2
				start := math.Atan2(-oy, -ox) - math.Pi/4
				arc := geom.Arc(0, 0, s/2, start, start+math.Pi/2, 16)
				pts := make([]core.Vec2, len(arc))
				for k, p := range arc {
					px, py := p.X+ox, p.Y+oy
					rx := px*math.Cos(ang) - py*math.Sin(ang)
					ry := px*math.Sin(ang) + py*math.Cos(ang)
					pts[k] = core.Vec2{X: cx + rx, Y: cy + ry}
				}
				scene.AddStroke(pts, false, lw, c, 1)
			}
		}
	}
	return scene, nil
}
