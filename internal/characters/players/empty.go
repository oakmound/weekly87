package players

import (
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
)

// EmptyConstructor for creating an empty space in a party select screen
// Used for party selection things
var EmptyConstructor *Constructor

// EmptyInit sets up the required placeholder renderables for empty constructor
// Most importantly sets runspeed to -1 for comparison checks in the future
func EmptyInit() {

	empty := render.NewEmptySprite(0, 0, 16, 32)

	emptyCharMap := map[string]render.Modifiable{
		"walkRT":    empty,
		"walkLT":    empty,
		"standRT":   empty,
		"standLT":   empty,
		"deadRT":    empty,
		"deadLT":    empty,
		"walkHold":  empty,
		"standHold": empty,
	}

	EmptyConstructor = &Constructor{
		AnimationMap: emptyCharMap,
		Dimensions:   floatgeom.Point2{16, 32},
		RunSpeed:     -1,
	}
}
