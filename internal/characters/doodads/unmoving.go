package doodads

import "github.com/oakmound/oak/physics"

type Unmoving struct{}

func (Unmoving) GetSpeed() physics.Vector {
	return physics.Vector{}
}

func (Unmoving) GetDelta() physics.Vector {
	return physics.Vector{}
}
