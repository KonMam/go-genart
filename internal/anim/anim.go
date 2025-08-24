package anim

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	stdpalette "image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"os"

	"genart/internal/config"
	"genart/internal/core"
	"genart/internal/palette"
	"genart/internal/render"
)

// Run executes an animated run based on cfg.Animation.
// It generates all frames, interpolates params and palette, and writes a GIF.
func Run(cfg *config.Config, eng core.Engine) error {
	anim := cfg.Animation
	if anim == nil {
		return fmt.Errorf("no animation section in config")
	}

	frames := int(math.Round(anim.Duration * float64(anim.FPS)))
	if frames <= 1 {
		return fmt.Errorf("invalid animation frames")
	}

	images := make([]*image.Paletted, 0, frames)
	delays := make([]int, 0, frames)

	for frame := 0; frame < frames; frame++ {
		tRaw := float64(frame) / float64(frames-1)
		t := ease(tRaw, anim.Easing)

		// Copy base params
		params := make(map[string]float64, len(cfg.Params))
		for k, v := range cfg.Params {
			params[k] = v
		}

		// Interpolate vary
		for key, raw := range anim.Vary {
			switch arr := raw.(type) {
			case []any:
				if len(arr) != 2 {
					continue
				}
				// scalar
				if f0, ok0 := arr[0].(float64); ok0 {
					if f1, ok1 := arr[1].(float64); ok1 {
						params[key] = lerp(f0, f1, t)
						continue
					}
				}
				// color
				if arr0, ok0 := toFloatSlice(arr[0]); ok0 && len(arr0) >= 3 {
					if arr1, ok1 := toFloatSlice(arr[1]); ok1 && len(arr1) >= 3 {
						c0 := core.RGBA{R: arr0[0], G: arr0[1], B: arr0[2], A: 1}
						c1 := core.RGBA{R: arr1[0], G: arr1[1], B: arr1[2], A: 1}
						cfg.Palette.Base = core.RGBA{
							R: lerp(c0.R, c1.R, t),
							G: lerp(c0.G, c1.G, t),
							B: lerp(c0.B, c1.B, t),
							A: lerp(c0.A, c1.A, t),
						}
					}
				}
			}
		}

		// rebuild palette
		var colors []core.RGBA
		switch cfg.Palette.Type {
		case "mono":
			colors = palette.Monochrome(cfg.Palette.Base, cfg.Palette.N)
		default:
			colors = palette.Monochrome(cfg.Palette.Base, cfg.Palette.N)
		}

		// derive seed
		subSeed := deriveSeed(cfg.Seed, eng.Name(), frame)
		rng := rand.New(rand.NewSource(subSeed))

		// generate
		scene, err := eng.Generate(context.Background(), rng, params, colors)
		if err != nil {
			return fmt.Errorf("engine failed: %w", err)
		}

		// render
		img, err := (render.GG{}).Render(scene, core.RenderConfig{
			Width:       cfg.Width,
			Height:      cfg.Height,
			Background:  cfg.Background,
			Margin:      cfg.Render.Margin,
			Supersample: cfg.Render.Supersample,
			Palette:     colors,
		})
		if err != nil {
			return fmt.Errorf("render failed: %w", err)
		}

		// convert to paletted
		pimg := image.NewPaletted(img.Bounds(), stdpalette.Plan9)
		draw.FloydSteinberg.Draw(pimg, img.Bounds(), img, image.Point{})
		images = append(images, pimg)
		delays = append(delays, int(100/anim.FPS))
	}

	// save GIF
	f, err := os.Create(cfg.Out)
	if err != nil {
		return fmt.Errorf("failed to create output: %w", err)
	}
	defer f.Close()

	return gif.EncodeAll(f, &gif.GIF{
		Image: images,
		Delay: delays,
	})
}

// --- helpers ---

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func toFloatSlice(v any) ([]float64, bool) {
	arr, ok := v.([]any)
	if !ok {
		return nil, false
	}
	out := make([]float64, len(arr))
	for i, x := range arr {
		f, ok := x.(float64)
		if !ok {
			return nil, false
		}
		out[i] = f
	}
	return out, true
}

func deriveSeed(root int64, label string, frame int) int64 {
	h := sha256.New()
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf[:8], uint64(root))
	binary.LittleEndian.PutUint64(buf[8:], uint64(frame))
	h.Write(buf)
	h.Write([]byte(label))
	sum := h.Sum(nil)
	return int64(binary.LittleEndian.Uint64(sum[:8]))
}

func ease(t float64, mode string) float64 {
	switch mode {
	case "cosine":
		// Smooth in/out
		return 0.5 - 0.5*math.Cos(t*math.Pi)
	case "sin":
		// Oscillating wave
		return 0.5 + 0.5*math.Sin(2*math.Pi*t)
	default:
		// Linear fallback
		return t
	}
}
