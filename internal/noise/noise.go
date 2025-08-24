package noise

// ScalarField2D is a contract for any scalar-valued 2D field.
type ScalarField2D interface {
	// At returns the field value at (x,y).
	// Typical contract: output in [-1,1].
	At(x, y float64) float64
}

// ScalarField3D is a contract for any scalar-valued 3D field.
type ScalarField3D interface {
	// At returns the field value at (x,y,z).
	// Typical contract: output in [-1,1].
	At(x, y, z float64) float64
}
