package noise

import (
	"math"

	"github.com/ojrac/opensimplex-go"
)

// LoopField is 2D simplex noise that also depends on a phase in [0,1) and
// is periodic in it: the phase walks a circle through the two extra
// dimensions of 4D noise, so phase 0 and phase 1 give the same field and a
// small phase step changes the field a little everywhere. It is the
// building block for parameters that must move a picture smoothly, and for
// animations that loop without a seam.
type LoopField struct {
	noise  opensimplex.Noise
	scale  float64
	radius float64
}

// NewLoopField returns a field whose features are about scale apart (in the
// same units as x and y). radius is how far one full phase cycle travels
// through noise space, in feature sizes: larger means more change per cycle.
func NewLoopField(seed int64, scale, radius float64) *LoopField {
	return &LoopField{noise: opensimplex.New(seed), scale: scale, radius: radius}
}

// At returns the field at (x, y) and phase, in about [-1, 1].
func (f *LoopField) At(x, y, phase float64) float64 {
	if f.scale <= 0 {
		return 0
	}
	a := 2 * math.Pi * phase
	return f.noise.Eval4(x/f.scale, y/f.scale, f.radius*math.Cos(a), f.radius*math.Sin(a))
}
