package noise

// Curl2D computes the 2D curl of a 2D scalar field.
// It returns a 2D vector (nx, ny) that is perpendicular to the gradient.
func Curl2D(field ScalarField2D, x, y, epsilon float64) (float64, float64) {
	n1 := field.At(x+epsilon, y)
	n2 := field.At(x-epsilon, y)
	n3 := field.At(x, y+epsilon)
	n4 := field.At(x, y-epsilon)

	dx := (n1 - n2) / (2 * epsilon)
	dy := (n3 - n4) / (2 * epsilon)

	return dy, -dx
}
