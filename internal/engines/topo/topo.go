// Package topo draws a contour map of a noise landscape: terraces filled by
// height and contour lines between them.
//
// All parameters are continuous. The landscape is periodic noise
// (noise.LoopField), so phase (0..1) morphs it and loops, and offset slides
// the contour levels uphill by one level per unit. Terrace colours repeat
// every 2(n-1) levels for an n-colour palette, so offset loops at 2(n-1)
// (6 for 4 colours), and at 1 with fill 0. The engine draws one random number whatever the parameters.
package topo

import (
	"context"
	"math"
	"math/rand"

	"genart/internal/core"
	"genart/internal/noise"
)

type Engine struct{}

func (Engine) Name() string { return "topo" }

func (Engine) Generate(_ context.Context, rng *rand.Rand, params map[string]float64, colors []core.RGBA) (core.Scene, error) {
	levels := core.Pick(params, "levels", 14)  // contour levels across the height range (continuous)
	freq := core.Pick(params, "freq", 1.6)     // hills across the picture
	phase := core.Pick(params, "phase", 0)     // landscape morph, 0..1 loops
	travel := core.Pick(params, "travel", 0.3) // how much one phase cycle changes the landscape
	offset := core.Pick(params, "offset", 0)   // contour levels slide by one level per unit, loops
	lw := core.Pick(params, "lw", 0.005)       // contour line width; 0 = none
	fill := core.Pick(params, "fill", 1)       // 0 = lines on background, 1 = terraces filled by height
	// style, when given, sets fill and lines together: 0 contour lines on
	// the background, 1 filled terraces with lines, 2 flat terraces.
	switch style := core.Pick(params, "style", -1); {
	case style < 0:
	case style < 0.5:
		fill = 0
	case style < 1.5:
		fill = 1
	default:
		fill, lw = 1, 0
	}
	res := int(core.Pick(params, "res", 240)) // marching-squares grid cells across

	land := noise.NewLoopField(rng.Int63(), 1, travel)
	height := func(x, y float64) float64 {
		// Zoom about the centre, so freq grows or shrinks the same hills.
		x, y = x-0.5, y-0.5
		h := land.At(x*freq, y*freq, phase) + 0.4*land.At(x*freq*2.1+11, y*freq*2.1, phase)
		// Simplex noise rarely leaves +-0.6; stretch it over 0..1.
		return 0.5 + 0.5*math.Tanh(2*h)
	}
	// level returns the continuous terrace coordinate: integers are contours.
	level := func(x, y float64) float64 { return height(x, y)*levels - offset }

	scene := core.Scene{}
	if fill > 0 {
		n := len(colors)
		scene.Items = append(scene.Items, core.Field{Color: func(x, y float64) core.RGBA {
			// A terrace keeps its colour as offset moves its contours,
			// so only pixels the contours pass over change. Colours
			// run up and back down the palette every 2(n-1) terraces,
			// which is why offset loops at that period.
			k := int(math.Floor(level(x, y)))
			p := 2 * (n - 1)
			k = ((k % p) + p) % p
			if k >= n {
				k = p - k
			}
			c := colors[k]
			c.A = fill
			return c
		}})
	}
	if lw > 0 {
		line := colors[0]
		if fill == 0 {
			line = colors[len(colors)-1]
		}
		for _, seg := range contours(level, res) {
			scene.AddStroke(seg[:], false, lw, line, 1)
		}
	}
	return scene, nil
}

// contours runs marching squares over [0,1]^2 on a res-wide grid and returns
// the segments where f crosses an integer.
func contours(f func(x, y float64) float64, res int) [][2]core.Vec2 {
	n := res + 1
	v := make([]float64, n*n)
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			v[j*n+i] = f(float64(i)/float64(res), float64(j)/float64(res))
		}
	}
	var segs [][2]core.Vec2
	h := 1 / float64(res)
	for j := 0; j < res; j++ {
		for i := 0; i < res; i++ {
			a, b, c, d := v[j*n+i], v[j*n+i+1], v[(j+1)*n+i+1], v[(j+1)*n+i]
			lo := math.Ceil(math.Min(math.Min(a, b), math.Min(c, d)))
			hi := math.Max(math.Max(a, b), math.Max(c, d))
			x0, y0 := float64(i)*h, float64(j)*h
			for t := lo; t <= hi; t++ {
				// Edge crossings in order: top, right, bottom, left.
				var pts []core.Vec2
				cross := func(p, q float64, px, py, qx, qy float64) {
					if (p < t) != (q < t) {
						s := (t - p) / (q - p)
						pts = append(pts, core.Vec2{X: x0 + (px+(qx-px)*s)*h, Y: y0 + (py+(qy-py)*s)*h})
					}
				}
				cross(a, b, 0, 0, 1, 0)
				cross(b, c, 1, 0, 1, 1)
				cross(d, c, 0, 1, 1, 1)
				cross(a, d, 0, 0, 0, 1)
				if len(pts) >= 2 {
					segs = append(segs, [2]core.Vec2{pts[0], pts[1]})
				}
				if len(pts) == 4 {
					segs = append(segs, [2]core.Vec2{pts[2], pts[3]})
				}
			}
		}
	}
	return segs
}
