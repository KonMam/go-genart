package anim

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	stdpalette "image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"genart/internal/config"
	"genart/internal/core"
	"genart/internal/palette"
	"genart/internal/render"
)

// FrameLog contains the parameters for a single frame of an animation.
type FrameLog struct {
	Frame  int                `json:"frame"`
	Params map[string]float64 `json:"params"`
}

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

	// An out without a .gif extension is a directory of numbered PNG
	// frames, full colour with no GIF quantisation.
	pngDir := filepath.Ext(cfg.Out) != ".gif"
	if pngDir {
		if err := os.MkdirAll(cfg.Out, 0o755); err != nil {
			return err
		}
	}

	images := make([]*image.Paletted, frames)
	delays := make([]int, frames)
	frameLogs := make([]FrameLog, frames)
	errs := make([]error, frames)

	// Frames are independent, so render them on every core.
	next := make(chan int)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for frame := range next {
				img, params, err := renderFrame(cfg, eng, frame, frames)
				if err != nil {
					errs[frame] = err
					continue
				}
				frameLogs[frame] = FrameLog{Frame: frame, Params: params}
				if pngDir {
					errs[frame] = writePNG(filepath.Join(cfg.Out, fmt.Sprintf("frame-%04d.png", frame)), img)
					continue
				}
				// convert to paletted
				pimg := image.NewPaletted(img.Bounds(), stdpalette.Plan9)
				draw.FloydSteinberg.Draw(pimg, img.Bounds(), img, image.Point{})
				images[frame] = pimg
				delays[frame] = int(100 / anim.FPS)
			}
		}()
	}
	for frame := range frames {
		next <- frame
	}
	close(next)
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	if anim.LogFrames {
		if err := logFrames(cfg.Out, frameLogs); err != nil {
			return fmt.Errorf("failed to log frames: %w", err)
		}
	}

	if pngDir {
		return nil
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

// renderFrame interpolates the params and palette for one frame and renders
// it.
func renderFrame(cfg *config.Config, eng core.Engine, frame, frames int) (image.Image, map[string]float64, error) {
	anim := cfg.Animation
	tRaw := float64(frame) / float64(frames-1)
	if anim.Loop {
		tRaw = float64(frame) / float64(frames)
	}
	t := ease(tRaw, anim.Easing)

	// Copy base params
	params := make(map[string]float64, len(cfg.Params))
	for k, v := range cfg.Params {
		params[k] = v
	}
	base := cfg.Palette.Base

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
					base = core.RGBA{
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
	colors := palette.Monochrome(base, cfg.Palette.N)

	// derive seed
	seedFrame := frame
	if anim.FixedSeed {
		seedFrame = 0
	}
	rng := rand.New(rand.NewSource(deriveSeed(cfg.Seed+config.VariantSeed(params), eng.Name(), seedFrame)))

	scene, err := eng.Generate(context.Background(), rng, params, colors)
	if err != nil {
		return nil, nil, fmt.Errorf("engine failed: %w", err)
	}
	img, err := (render.GG{}).Render(scene, core.RenderConfig{
		Width:       cfg.Width,
		Height:      cfg.Height,
		Background:  cfg.Background,
		Margin:      cfg.Render.Margin,
		Supersample: cfg.Render.Supersample,
		Palette:     colors,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("render failed: %w", err)
	}
	return img, params, nil
}

func writePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func logFrames(out string, logs []FrameLog) error {
	ext := filepath.Ext(out)
	logFile := out[0:len(out)-len(ext)] + ".json"

	f, err := os.Create(logFile)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(logs)
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
