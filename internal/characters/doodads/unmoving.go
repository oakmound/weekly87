package doodads

import "github.com/oakmound/oak/physics"

// Unmoving objects still need to return that they have a speed and delta.
// The speed and delta are always zero values though!
type Unmoving struct{}

// GetSpeed returns the current speed (always zero value)
func (Unmoving) GetSpeed() physics.Vector {
	return physics.Vector{}
}

// GetDelta returns the current delta (always zero value)
func (Unmoving) GetDelta() physics.Vector {
	return physics.Vector{}
}
