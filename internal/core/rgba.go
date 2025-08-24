package core

import (
	"encoding/json"
	"fmt"
)

// A color with float components in [0,1].
type RGBA struct {
	R, G, B, A float64
}

// UnmarshalJSON allows RGBA to be decoded from either
// an array [r,g,b,a] or an object {"R":..,"G":..,"B":..,"A":..}.
func (c *RGBA) UnmarshalJSON(data []byte) error {
	// Try array form first
	var arr []float64
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) < 3 || len(arr) > 4 {
			return fmt.Errorf("invalid RGBA array length: %d", len(arr))
		}
		c.R, c.G, c.B = arr[0], arr[1], arr[2]
		if len(arr) == 4 {
			c.A = arr[3]
		} else {
			c.A = 1.0
		}
		return nil
	}

	// Try object form
	var obj struct {
		R float64 `json:"R"`
		G float64 `json:"G"`
		B float64 `json:"B"`
		A float64 `json:"A"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	c.R, c.G, c.B, c.A = obj.R, obj.G, obj.B, obj.A
	return nil
}

// MarshalJSON always encodes RGBA as a 4-element array [r,g,b,a].
func (c RGBA) MarshalJSON() ([]byte, error) {
	arr := [4]float64{c.R, c.G, c.B, c.A}
	return json.Marshal(arr)
}
