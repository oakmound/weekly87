package characters

import (
	"path/filepath"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var SpearmanConstructor *PlayerConstructor

func Init() {
	animFilePath := (filepath.Join("16x32", "warrior.png"))
	sheet, err := render.LoadSprites(filepath.Join("assets", "images"),
		animFilePath, 16, 32, 0)
	dlog.ErrorCheck(err)
	standRT := sheet[0][0]
	standLT := sheet[0][0].Copy().Modify(mod.FlipX)

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	SpearmanConstructor = &PlayerConstructor{
		AnimationMap: map[string]render.Modifiable{
			"walkRT":  walkRT,
			"walkLT":  walkLT,
			"standRT": standRT,
			"standLT": standLT,
		},
		Dimensions: floatgeom.Point2{16, 32},
		Speed:      floatgeom.Point2{0, 5},
		RunSpeed:   6.0,
	}
}

func NewSpearman(x, y float64) (*Player, error) {
	cs := SpearmanConstructor.Copy()
	cs.Position = floatgeom.Point2{x, y}
	return cs.NewPlayer()
}
