package players

import (
	"path/filepath"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var SpearmanConstructor *Constructor

func Init() {
	animFilePath := (filepath.Join("16x32", "warrior.png"))
	sheet, err := render.LoadSprites(filepath.Join("assets", "images"),
		animFilePath, 16, 32, 0)
	dlog.ErrorCheck(err)
	standRT := sheet[0][0].Copy()
	standLT := sheet[0][0].Copy().Modify(mod.FlipX)
	standHold := sheet[0][1].Copy().Modify(mod.FlipX)

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	walkHold, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 1, 2, 1, 0, 1}...)
	dlog.ErrorCheck(err)
	walkHold = walkHold.Copy().Modify(mod.FlipX).(*render.Sequence)

	deadRT := sheet[0][0].Copy()
	deadRT.Filter(mod.Fade(125))
	deadRT.Filter(mod.ColorBalance(-50, -50, -50))

	deadLT := deadRT.Copy().Modify(mod.FlipX)

	SpearmanConstructor = &Constructor{
		AnimationMap: map[string]render.Modifiable{
			"walkRT":    walkRT,
			"walkLT":    walkLT,
			"standRT":   standRT,
			"standLT":   standLT,
			"deadRT":    deadRT,
			"deadLT":    deadLT,
			"walkHold":  walkHold,
			"standHold": standHold,
		},
		Dimensions: floatgeom.Point2{16, 32},
		Speed:      floatgeom.Point2{0, 5},
		RunSpeed:   3.0,
	}
}

func NewSpearman(x, y float64) (*Player, error) {
	cs := SpearmanConstructor.Copy()
	cs.Position = floatgeom.Point2{x, y}
	return cs.NewPlayer()
}
