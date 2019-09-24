package players

import (
	"path/filepath"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var InnKeeperConstructor *Constructor

func InnKeeperInit() {

	sh, err := render.LoadSheet(filepath.Join("assets", "images", "16x32"), "innkeeper.png", 16, 32, 0)
	dlog.ErrorCheck(err)

	sheet := sh.ToSprites()

	//Spearman
	standRT := sheet[0][0].Copy()
	standLT := sheet[0][0].Copy().Modify(mod.FlipX)
	standHold := sheet[0][1].Copy().Modify(mod.FlipX)

	walkRT, err := render.NewSheetSequence(sh, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)

	walkLT := walkRT.Copy().Modify(mod.FlipX)

	walkHold, err := render.NewSheetSequence(sh, 8, []int{1, 1, 2, 1, 0, 1}...)
	dlog.ErrorCheck(err)
	walkHold = walkHold.Copy().Modify(mod.FlipX).(*render.Sequence)

	animMap := map[string]render.Modifiable{
		"walkRT":    walkRT,
		"walkLT":    walkLT,
		"standRT":   standRT,
		"standLT":   standLT,
		"deadRT":    walkRT,
		"deadLT":    walkLT,
		"walkHold":  walkHold,
		"standHold": standHold,
	}

	InnKeeperConstructor = &Constructor{
		AnimationMap: animMap,
		Dimensions:   floatgeom.Point2{16, 32},
	}
}
