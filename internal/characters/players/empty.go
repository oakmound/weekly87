package players

import (
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
)

var EmptyConstructor *Constructor

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
