package config

import "genart/internal/core"

// Config represents a full run configuration.
// This is both the INPUT (when loaded from file/string)
// and the OUTPUT (what we dump after a run).
type Config struct {
	Engine     string             `json:"engine"`
	Width      int                `json:"width"`
	Height     int                `json:"height"`
	Out        string             `json:"out"`
	Seed       int64              `json:"seed"`
	Background core.RGBA          `json:"bg"`
	Palette    PaletteConfig      `json:"palette"`
	Params     map[string]float64 `json:"params"`

	Render    RenderConfig     `json:"render"`
	Animation *AnimationConfig `json:"animation,omitempty"`
}

// PaletteConfig controls palette generation.
type PaletteConfig struct {
	Type string    `json:"type"` // e.g. "mono"
	Base core.RGBA `json:"base"` // base color
	N    int       `json:"n"`    // number of colors
}

// RenderConfig controls renderer settings.
type RenderConfig struct {
	Margin      float64 `json:"margin"`
	Supersample int     `json:"supersample"`
}

// AnimationConfig controls animation runs.
// If nil, the run is static (PNG).
type AnimationConfig struct {
	Duration float64        `json:"duration"` // seconds
	FPS      int            `json:"fps"`
	Vary     map[string]any `json:"vary,omitempty"`   // param name -> [start,end]
	Easing   string         `json:"easing,omitempty"` // "linear" (default), "cosine", "sin"
	LogFrames bool          `json:"log_frames,omitempty"`
}
