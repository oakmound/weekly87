package players

import (
	"path/filepath"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/recolor"
)

var MageConstructor *Constructor
var WhiteMageConstructor *Constructor

func MageInit() {
	animFilePath := (filepath.Join("16x32", "mage.png"))
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

	ghostFilePath := filepath.Join("16x32", "mageghost.png")
	dlog.ErrorCheck(err)

	deadRT, err := render.LoadSheetSequence(ghostFilePath, 16, 32, 0, 8, []int{0, 0, 1, 0}...)
	dlog.ErrorCheck(err)
	deadLT := deadRT.Copy().Modify(mod.FlipX)

	mageCharMap := map[string]render.Modifiable{
		"walkRT":    walkRT,
		"walkLT":    walkLT,
		"standRT":   standRT,
		"standLT":   standLT,
		"deadRT":    deadRT,
		"deadLT":    deadLT,
		"walkHold":  walkHold,
		"standHold": standHold,
	}

	MageConstructor = &Constructor{
		AnimationMap: mageCharMap,
		Dimensions:   floatgeom.Point2{16, 32},
		Speed:        floatgeom.Point2{0, 5},
		RunSpeed:     3.0,
		Special1:     abilities.Fireball,
		Special2:     abilities.Invulnerability,
	}

	whiteMageMap := filterCharMap(mageCharMap, recolor.Recolor(recolor.WhiteMage))

	WhiteMageConstructor = &Constructor{
		AnimationMap: whiteMageMap,
		Dimensions:   floatgeom.Point2{16, 32},
		Speed:        floatgeom.Point2{0, 5},
		RunSpeed:     3.0,
		Special1:     abilities.Fireball,
		Special2:     abilities.Fireball,
	}

}
