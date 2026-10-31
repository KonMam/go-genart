package render

import (
	"errors"
	"image"
	"math"

	"genart/internal/core"

	"github.com/fogleman/gg"
)

type GG struct{}

func (GG) Name() string { return "gg" }

// Render maps logical [0..1] coordinates to pixels, applies margins,
// and paints fills first, then strokes.
func (GG) Render(scene core.Scene, cfg core.RenderConfig) (image.Image, error) {
	W, H := cfg.Width, cfg.Height
	if W <= 0 || H <= 0 {
		return nil, ErrInvalidSize
	}

	dc := gg.NewContext(W, H)

	// Background
	dc.SetRGBA(cfg.Background.R, cfg.Background.G, cfg.Background.B, cfg.Background.A)
	dc.Clear()

	// margin in pixels
	marginPx := cfg.Margin * float64(min(W, H))
	x0 := marginPx
	y0 := marginPx
	x1 := float64(W) - marginPx
	y1 := float64(H) - marginPx
	sx := x1 - x0
	sy := y1 - y0
	minWH := float64(min(W, H))

	mapPt := func(v core.Vec2) (float64, float64) {
		return x0 + v.X*sx, y0 + v.Y*sy
	}

	// Fields first, pixel by pixel
	for _, it := range scene.Items {
		f, ok := it.(core.Field)
		if !ok {
			continue
		}
		img := dc.Image().(*image.RGBA)
		for py := int(y0); py < int(y1+0.5); py++ {
			for px := int(x0); px < int(x1+0.5); px++ {
				c := f.Color((float64(px)+0.5-x0)/sx, (float64(py)+0.5-y0)/sy)
				o := img.PixOffset(px, py)
				blend := func(dst uint8, v float64) uint8 {
					v = math.Max(0, math.Min(1, v))
					return uint8(math.Round(float64(dst)*(1-c.A) + 255*v*c.A))
				}
				img.Pix[o] = blend(img.Pix[o], c.R)
				img.Pix[o+1] = blend(img.Pix[o+1], c.G)
				img.Pix[o+2] = blend(img.Pix[o+2], c.B)
			}
		}
	}

	// Render items
	for _, it := range scene.Items {
		switch s := it.(type) {
		case core.Fill:
			dc.NewSubPath()
			for i, p := range s.Polygon.Points {
				x, y := mapPt(p)
				if i == 0 {
					dc.MoveTo(x, y)
				} else {
					dc.LineTo(x, y)
				}
			}
			if s.Polygon.Closed {
				dc.ClosePath()
			}
			dc.SetRGBA(s.Color.R, s.Color.G, s.Color.B, s.Alpha)
			dc.Fill()

		case core.Stroke:
			dc.NewSubPath()
			for i, p := range s.Path.Points {
				x, y := mapPt(p)
				if i == 0 {
					dc.MoveTo(x, y)
				} else {
					dc.LineTo(x, y)
				}
			}
			if s.Path.Closed {
				dc.ClosePath()
			}
			dc.SetRGBA(s.Color.R, s.Color.G, s.Color.B, s.Alpha)
			dc.SetLineWidth(s.Width * minWH) // logical → pixels
			dc.Stroke()
		}
	}

	return dc.Image(), nil
}

// simple helper
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// custom error
var ErrInvalidSize = errors.New("render: width and height must be > 0")
