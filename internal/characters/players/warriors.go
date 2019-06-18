package players

import (
	"image/color"
	"path/filepath"

	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/recolor"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var WarriorConstructors = map[string]*Constructor{}

func WarriorsInit() {

	var warriorDefinitions = []ClassDefinition{
		{
			Name: "Paladin",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{240, 240, 240, 70},
			},
			Special1: abilities.SwordSwipe,
			Special2: abilities.Shield,
		},
		{
			Name: "Berserker",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{100, 100, 240, 20},
			},
			Special1: abilities.SwordSwipe,
			Special2: abilities.Shield,
		},
		{
			Name: "Spearman",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{100, 240, 100, 150},
			},
			Special1: abilities.SwordSwipe,
			Special2: abilities.Shield,
		},
		{
			Name: "Swordsman",
			//LayerColors: map[string]color.RGBA{
			//	"clothes": color.RGBA{240, 100, 100, 125},
			//},
			Special1: abilities.SwordSwipe,
			Special2: abilities.SwordSwipe,
		},
	}

	for _, def := range warriorDefinitions {

		animFilePath := (filepath.Join("16x32", "warrior.png"))
		sheet, err := render.LoadSprites(filepath.Join("assets", "images"),
			animFilePath, 16, 32, 0)
		dlog.ErrorCheck(err)

		ghostFilePath := filepath.Join("16x32", "warriorghost.png")
		dlog.ErrorCheck(err)

		deadRT, err := render.LoadSheetSequence(ghostFilePath, 16, 32, 0, 8, []int{0, 0, 1, 0}...)
		dlog.ErrorCheck(err)
		deadLT := deadRT.Copy().Modify(mod.FlipX)

		//Spearman
		standRT := sheet[0][0].Copy()
		standLT := sheet[0][0].Copy().Modify(mod.FlipX)
		standHold := sheet[0][1].Copy().Modify(mod.FlipX)

		walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
		dlog.ErrorCheck(err)

		walkLT := walkRT.Copy().Modify(mod.FlipX)

		walkHold, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 1, 2, 1, 0, 1}...)
		dlog.ErrorCheck(err)
		walkHold = walkHold.Copy().Modify(mod.FlipX).(*render.Sequence)

		warriorMap := map[string]render.Modifiable{
			"walkRT":    walkRT,
			"walkLT":    walkLT,
			"standRT":   standRT,
			"standLT":   standLT,
			"deadRT":    deadRT,
			"deadLT":    deadLT,
			"walkHold":  walkHold,
			"standHold": standHold,
		}

		if def.Name == "Swordsman" {
			filterCharMap(warriorMap, recolor.Recolor(recolor.WarriorSwordsman))
		}

		WarriorConstructors[def.Name] = &Constructor{
			AnimationMap: warriorMap,
			Dimensions:   floatgeom.Point2{16, 32},
			Speed:        floatgeom.Point2{0, 5},
			RunSpeed:     3.0,
			Special1:     def.Special1,
			Special2:     def.Special2,
		}

	}

}
