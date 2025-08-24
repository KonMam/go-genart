package noise

import "math"

// Remap01 maps noise in [-1,1] to [0,1].
func Remap01(v float64) float64 {
	return 0.5 * (v + 1)
}

// RemapRange maps noise in [-1,1] into [a,b].
func RemapRange(v, a, b float64) float64 {
	return a + Remap01(v)*(b-a)
}

// ToAngle maps noise in [-1,1] to [0,2π).
func ToAngle(v float64) float64 {
	return (v + 1) * math.Pi
}
