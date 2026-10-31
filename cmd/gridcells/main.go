// Command gridcells renders the cells of an animated parameter grid for
// gbemu's gbart grid: every combination of the grid's param values is one
// cell, a seamless loop rendered from the config's animation, written as
//
//	<out>/cell-NNNN/frame-KK.png
//
// where NNNN is the combination's index with the last param varying
// fastest. -grid is the grid.json that gbart reads (name, fps, params with
// values or from/to/steps).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"genart/internal/anim"
	"genart/internal/config"
	"genart/internal/core"
	"genart/internal/engines/registry"
)

type gridParam struct {
	Name   string    `json:"name"`
	Values []float64 `json:"values"`
	From   *float64  `json:"from"`
	To     *float64  `json:"to"`
	Steps  int       `json:"steps"`
}

type grid struct {
	Name   string      `json:"name"`
	FPS    int         `json:"fps"`
	Params []gridParam `json:"params"`
}

func main() {
	cfgPath := flag.String("config", "", "go-genart config with an animation section (required)")
	gridPath := flag.String("grid", "", "grid.json (required)")
	out := flag.String("out", "", "cells directory (required)")
	width := flag.Int("width", 320, "cell width")
	height := flag.Int("height", 288, "cell height")
	flag.Parse()
	if *cfgPath == "" || *gridPath == "" || *out == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(*cfgPath, *gridPath, *out, *width, *height); err != nil {
		fmt.Fprintln(os.Stderr, "gridcells:", err)
		os.Exit(1)
	}
}

func run(cfgPath, gridPath, out string, width, height int) error {
	base, err := config.Load(cfgPath)
	if err != nil {
		return err
	}
	if base.Animation == nil {
		return fmt.Errorf("%s has no animation section", cfgPath)
	}
	eng, ok := registry.Engines()[base.Engine]
	if !ok {
		return fmt.Errorf("unknown engine %q", base.Engine)
	}
	data, err := os.ReadFile(gridPath)
	if err != nil {
		return err
	}
	var g grid
	if err := json.Unmarshal(data, &g); err != nil {
		return fmt.Errorf("%s: %w", gridPath, err)
	}
	counts := make([]int, len(g.Params))
	cells := 1
	for i := range g.Params {
		p := &g.Params[i]
		if len(p.Values) == 0 {
			if p.From == nil || p.To == nil || p.Steps < 1 {
				return fmt.Errorf("param %q needs values or from, to and steps", p.Name)
			}
			for k := range p.Steps {
				t := 0.0
				if p.Steps > 1 {
					t = float64(k) / float64(p.Steps-1)
				}
				p.Values = append(p.Values, *p.From+(*p.To-*p.From)*t)
			}
		}
		counts[i] = len(p.Values)
		cells *= counts[i]
	}
	frames := int(math.Round(base.Animation.Duration * float64(base.Animation.FPS)))

	type job struct{ cell, frame int }
	jobs := make(chan job)
	errs := make(chan error, 1)
	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if err := renderCell(base, eng, g, counts, j.cell, j.frame, frames, out, width, height); err != nil {
					select {
					case errs <- err:
					default:
					}
				}
			}
		}()
	}
	for c := range cells {
		if err := os.MkdirAll(filepath.Join(out, fmt.Sprintf("cell-%04d", c)), 0o755); err != nil {
			return err
		}
		for f := range frames {
			jobs <- job{c, f}
		}
	}
	close(jobs)
	wg.Wait()
	select {
	case err := <-errs:
		return err
	default:
	}
	fmt.Printf("%s: %d cells x %d frames -> %s\n", g.Name, cells, frames, out)
	return nil
}

// renderCell renders one frame of one cell: the base config with the
// cell's param values, at the cell size.
func renderCell(base *config.Config, eng core.Engine, g grid, counts []int, cell, frame, frames int, out string, w, h int) error {
	cfg := *base
	cfg.Width, cfg.Height = w, h
	cfg.Params = make(map[string]float64, len(base.Params)+len(g.Params))
	for k, v := range base.Params {
		cfg.Params[k] = v
	}
	idx := cell
	for i := len(g.Params) - 1; i >= 0; i-- {
		cfg.Params[g.Params[i].Name] = g.Params[i].Values[idx%counts[i]]
		idx /= counts[i]
	}
	img, _, err := anim.RenderFrame(&cfg, eng, frame, frames)
	if err != nil {
		return fmt.Errorf("cell %d frame %d: %w", cell, frame, err)
	}
	f, err := os.Create(filepath.Join(out, fmt.Sprintf("cell-%04d", cell), fmt.Sprintf("frame-%02d.png", frame)))
	if err != nil {
		return err
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
