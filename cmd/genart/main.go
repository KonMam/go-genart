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

	"genart/internal/config"
	"genart/internal/core"
	"genart/internal/engines/blackhole"
	"genart/internal/engines/contourlines"
	"genart/internal/engines/flowfield"
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
	}

	eng, ok := engines[cfg.Engine]
	if !ok {
		exitErr(fmt.Sprintf("invalid engine %q", cfg.Engine))
	}

	// --- Validate ---
	if cfg.Width <= 0 || cfg.Height <= 0 {
		exitErr("width and height must be > 0")
	}
	if cfg.Out == "" {
		exitErr("output path cannot be empty")
	}

	// --- Build palette ---
	var colors []core.RGBA
	switch cfg.Palette.Type {
	case "mono":
		colors = palette.Monochrome(cfg.Palette.Base, cfg.Palette.N)
	default:
		exitErr(fmt.Sprintf("unknown palette type %q", cfg.Palette.Type))
	}

	// --- Print root seed ---
	fmt.Fprintf(os.Stderr, "Root seed: %d\n", cfg.Seed)

	// --- Derive per-engine seed ---
	subSeed := deriveSeed(cfg.Seed, eng.Name())
	rng := rand.New(rand.NewSource(subSeed))

	// --- Engine Generate ---
	scene, err := eng.Generate(context.Background(), rng, cfg.Params, colors)
	if err != nil {
		exitErr("engine failed: " + err.Error())
	}

	// --- Render ---
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

	// --- Save PNG ---
	f, err := os.Create(cfg.Out)
	if err != nil {
		exitErr("failed to create file: " + err.Error())
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		exitErr("failed to encode PNG: " + err.Error())
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
