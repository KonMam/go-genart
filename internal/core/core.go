package core

import (
	"context"
	"image"
	"math/rand"
)

// A 2D point in logical coordinates [0,1] × [0,1].
type Vec2 struct {
	X, Y float64
}

// A sequence of points. Closed indicates whether
// the path should loop back to the first point.
type Path struct {
	Points []Vec2
	Closed bool
}


// A marker interface for things that can be drawn in a Scene.
type Item interface{ isItem() }

// Stroke draws an outline along a path.
type Stroke struct {
	Path  Path
	Width float64 // logical fraction of min(width,height)
	Color RGBA
	Alpha float64
}

func (Stroke) isItem() {}

// Fill draws a filled polygon.
type Fill struct {
	Polygon Path
	Color   RGBA
	Alpha   float64
}

func (Fill) isItem() {}

// A collection of drawing instructions.
type Scene struct {
	Items []Item
}

// Generates a Scene from parameters and randomness.
type Engine interface {
	Name() string
	Generate(ctx context.Context, rng *rand.Rand, params map[string]float64, colors []RGBA) (Scene, error)
}

// Controls how a Scene is mapped into pixels.
type RenderConfig struct {
	Width, Height int
	Background    RGBA
	Margin        float64 // fraction of min(width,height)
	Supersample   int
	Palette       []RGBA
}

// Paints a Scene to an image.
type Renderer interface {
	Name() string
	Render(scene Scene, cfg RenderConfig) (image.Image, error)
}
