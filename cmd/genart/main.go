package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Engine string             `json:"engine"`
	Width  int                `json:"width"`
	Height int                `json:"height"`
	Out    string             `json:"out"`
	Seed   int64              `json:"seed"`
	Params map[string]float64 `json:"params"`
}

var allowedEngines = map[string]bool{
	"square": true,
	"circle": true,
}

func main() {
	engine := flag.String("engine", "square", "engine to use (square|circle)")
	width := flag.Int("w", 1000, "output width in pixels")
	height := flag.Int("h", 1000, "output height in pixels")
	out := flag.String("out", "out.png", "output PNG path")
	seed := flag.Int64("seed", 42, "random seed")
	paramsCSV := flag.String("params", "", "engine params as k=v,k=v (numbers)")
	flag.Parse()

	if !allowedEngines[*engine] {
		exitErr(fmt.Sprintf("invalid engine %q (allowed: square,circle)", *engine))
	}
	if *width <= 0 {
		exitErr("width must be > 0")
	}
	if *height <= 0 {
		exitErr("height must be > 0")
	}
	if *out == "" {
		exitErr("output path cannot be empty")
	}

	params, err := parseParams(*paramsCSV)
	if err != nil {
		exitErr(err.Error())
	}

	cfg := Config{
		Engine: *engine,
		Width:  *width,
		Height: *height,
		Out:    *out,
		Seed:   *seed,
		Params: params,
	}

	// Print config as JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		exitErr("failed to encode config: " + err.Error())
	}
}

// parseParams parses "k=v,k=v" into map[string]float64
func parseParams(csv string) (map[string]float64, error) {
	m := make(map[string]float64)
	if csv == "" {
		return m, nil
	}
	for part := range strings.SplitSeq(csv, ",") {
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
