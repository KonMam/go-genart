package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"math/rand"
	"os"

	"genart/internal/anim"
	"genart/internal/config"
	"genart/internal/core"
	"genart/internal/engines/blackhole"
	"genart/internal/engines/contourlines"
	"genart/internal/engines/flowfield"
	"genart/internal/engines/perlinpearls"
	"genart/internal/engines/swirl"
	"genart/internal/engines/flow"
	"genart/internal/engines/strata"
	"genart/internal/palette"
	"genart/internal/render"
)

func main() {
	// --- Flags ---
	configFlag := flag.String("config", "", "JSON config string or path to .json file")
	flag.Parse()

	if *configFlag == "" {
		exitErr("you must pass -config (JSON string or file)")
	}

	// --- Load config ---
	cfg, err := config.Load(*configFlag)
	if err != nil {
		exitErr("failed to load config: " + err.Error())
	}

	// --- Registry ---
	engines := map[string]core.Engine{
		"flowfield":    flowfield.Engine{},
		"contourlines": contourlines.Engine{},
		"blackhole":    blackhole.Engine{},
		"perlinpearls": perlinpearls.Engine{},
		"swirl":        swirl.Engine{},
		"flow":         flow.Engine{},
		"strata":       strata.Engine{},
	}

	eng, ok := engines[cfg.Engine]
	if !ok {
		exitErr(fmt.Sprintf("invalid engine %q", cfg.Engine))
	}

	// --- Build palette ---
	var colors []core.RGBA
	switch cfg.Palette.Type {
	case "mono":
		colors = palette.Monochrome(cfg.Palette.Base, cfg.Palette.N)
	case "split-complementary":
		colors = palette.SplitComplementary(cfg.Palette.Base, cfg.Palette.N)
	case "analogous":
		colors = palette.Analogous(cfg.Palette.Base, cfg.Palette.N)
	default:
		exitErr(fmt.Sprintf("unknown palette type %q", cfg.Palette.Type))
	}

	// --- Print root seed ---
	fmt.Fprintf(os.Stderr, "Root seed: %d\n", cfg.Seed)

	// --- Run animation or static render ---
	if cfg.Animation != nil {
		if err := anim.Run(cfg, eng); err != nil {
			exitErr("animation failed: " + err.Error())
		}
	} else {
		// Static run
		subSeed := deriveSeed(cfg.Seed, eng.Name())
		rng := rand.New(rand.NewSource(subSeed))

		scene, err := eng.Generate(context.Background(), rng, cfg.Params, colors)
		if err != nil {
			exitErr("engine failed: " + err.Error())
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
			exitErr("render failed: " + err.Error())
		}

		f, err := os.Create(cfg.Out)
		if err != nil {
			exitErr("failed to create file: " + err.Error())
		}
		defer f.Close()
		if err := png.Encode(f, img); err != nil {
			exitErr("failed to encode PNG: " + err.Error())
		}
	}

	// --- Print final config JSON ---
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(cfg)
}

// --- Helpers ---

func deriveSeed(root int64, label string) int64 {
	h := sha256.New()
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(root))
	h.Write(buf)
	h.Write([]byte(label))
	sum := h.Sum(nil)
	return int64(binary.LittleEndian.Uint64(sum[:8]))
}

func exitErr(msg string) {
	fmt.Fprintln(os.Stderr, "genart:", msg)
	os.Exit(2)
}
