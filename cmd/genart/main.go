package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"genart/internal/core"
	"genart/internal/engines/circle"
	"genart/internal/engines/square"
	"genart/internal/render"
)

func main() {
	// --- Flags ---
	engine := flag.String("engine", "square", "engine to use (square|circle)")
	width := flag.Int("w", 1000, "output width in pixels")
	height := flag.Int("h", 1000, "output height in pixels")
	out := flag.String("out", "out.png", "output PNG path")
	seed := flag.Int64("seed", 42, "random seed")
	paramsCSV := flag.String("params", "", "engine params as k=v,k=v (numbers)")
	flag.Parse()

	// --- Registry ---
	engines := map[string]core.Engine{
		"square": square.Engine{},
		"circle": circle.Engine{},
	}

	eng, ok := engines[*engine]
	if !ok {
		exitErr(fmt.Sprintf("invalid engine %q (allowed: square|circle)", *engine))
	}

	// --- Validate ---
	if *width <= 0 || *height <= 0 {
		exitErr("width and height must be > 0")
	}
	if *out == "" {
		exitErr("output path cannot be empty")
	}

	// --- Parse params ---
	params, err := parseParams(*paramsCSV)
	if err != nil {
		exitErr(err.Error())
	}

	// --- Engine Generate ---
	rng := rand.New(rand.NewSource(*seed)) // currently ignored, but passed in
	scene, err := eng.Generate(context.Background(), rng, params)
	if err != nil {
		exitErr("engine failed: " + err.Error())
	}

	// --- Render ---
	img, err := (render.GG{}).Render(scene, core.RenderConfig{
		Width:      *width,
		Height:     *height,
		Background: core.RGBA{R: 1,G: 1,B: 1,A: 1},
		Margin:     0.05,
		Supersample: 1,
	})
	if err != nil {
		exitErr("render failed: " + err.Error())
	}

	// --- Save PNG ---
	f, err := os.Create(*out)
	if err != nil {
		exitErr("failed to create file: " + err.Error())
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		exitErr("failed to encode PNG: " + err.Error())
	}

	// --- Print summary JSON ---
	summary := map[string]any{
		"engine": eng.Name(),
		"width":  *width,
		"height": *height,
		"out":    *out,
		"seed":   *seed,
		"params": params,
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(summary)
}

// --- Helpers ---

func parseParams(csv string) (map[string]float64, error) {
	m := make(map[string]float64)
	if csv == "" {
		return m, nil
	}
	for _, part := range strings.Split(csv, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid param %q (expected k=v)", part)
		}
		key := strings.TrimSpace(kv[0])
		valStr := strings.TrimSpace(kv[1])
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value for param %q: %q (want number)", key, valStr)
		}
		m[key] = val
	}
	return m, nil
}

func exitErr(msg string) {
	fmt.Fprintln(os.Stderr, "genart:", msg)
	os.Exit(2)
}
