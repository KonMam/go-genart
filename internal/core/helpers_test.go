package core

import "testing"

func TestNewStrokeAndFill(t *testing.T) {
	scene := Scene{}

	pts := []Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}
	color := RGBA{R: 1, G: 0, B: 0, A: 1}

	scene.AddStroke(pts, false, 0.01, color, 0.8)
	scene.AddFill(pts, color, 0.5)

	if len(scene.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(scene.Items))
	}

	if stroke, ok := scene.Items[0].(Stroke); !ok {
		t.Fatalf("expected Stroke, got %T", scene.Items[0])
	} else if stroke.Path.Closed {
		t.Errorf("expected open path, got closed")
	}

	if fill, ok := scene.Items[1].(Fill); !ok {
		t.Fatalf("expected Fill, got %T", scene.Items[1])
	} else if !fill.Polygon.Closed {
		t.Errorf("expected closed polygon, got open")
	}
}
