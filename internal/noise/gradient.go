package noise

// Gradient2D computes the gradient (∇f) of a 2D noise field at (x,y).
// Uses central finite differences with epsilon.
func Gradient2D(f ScalarField2D, x, y, eps float64) (dx, dy float64) {
	dx = (f.At(x+eps, y) - f.At(x-eps, y)) / (2 * eps)
	dy = (f.At(x, y+eps) - f.At(x, y-eps)) / (2 * eps)
	return
}
