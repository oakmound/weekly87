package players

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/recolor"
	"github.com/solovev/gopsd"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var WarriorConstructors = map[string]*Constructor{}

func WarriorsInit() {

	var warriorDefinitions = []ClassDefinition{
		{
			Name: "Swordsman",
			//LayerColors: map[string]color.RGBA{
			//	"clothes": color.RGBA{240, 100, 100, 125},
			//},
			Special1: abilities.SwordSwipe,
			Special2: abilities.SelfShield,
		},
		{
			Name: "Paladin",
			LayerColors: map[string]color.RGBA{
				//"clothes": color.RGBA{240, 240, 240, 70},
				"clothes": color.RGBA{200, 200, 200, 120},
			},
			Special1: abilities.HammerSmack,
			Special2: abilities.PartyShield,
		},
		{
			Name: "Berserker",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{160, 70, 70, 90},
			},
			Special1: abilities.SwordSwipe,
			Special2: abilities.Rage,
		},
		{
			Name: "Spearman",
			LayerColors: map[string]color.RGBA{
				"clothes": color.RGBA{70, 70, 150, 90},
			},
			Special1: abilities.SpearStab,
			Special2: abilities.SpearThrow,
		},
	}

	for _, def := range warriorDefinitions {
		psdFilePath := filepath.Join("assets", "images", "16x32", "warrior.psd")
		psd, err := gopsd.ParseFromPath(psdFilePath)

		combined := render.NewCompositeM()

		for _, layer := range psd.Layers {
			//TODO: combine strat with that of mage
			img, err := layer.GetImage()
			dlog.ErrorCheck(err)
			rgba, ok := img.(*image.RGBA)
			if !ok {
				dlog.Error("Image was not RGBA in underlying type")
			}
			sp := render.NewSprite(float64(layer.Rectangle.X), float64(layer.Rectangle.Y), rgba)
			if c, ok := def.LayerColors[strings.ToLower(layer.Name)]; ok {
				fmt.Println("We found the right layer", layer.Name)
				// Recolor this layer
				sp.Filter(recolor.WithStrategy(recolor.ColorMix(c)))
			}
			// Add this layer to the combined image
			// Todo: bug with shoulder having some pixel flashing
			combined.Append(sp)
		}

		// flatten composite
		combinedSp := combined.ToSprite()
		sh, err := render.MakeSheet(combinedSp.GetRGBA(), 16, 32, 0)
		dlog.ErrorCheck(err)
		sheet := sh.ToSprites()

		ghostFilePath := filepath.Join("16x32", "warriorghost.png")
		dlog.ErrorCheck(err)

		deadRT, err := render.LoadSheetSequence(ghostFilePath, 16, 32, 0, 8, []int{0, 0, 1, 0}...)
		dlog.ErrorCheck(err)
		deadLT := deadRT.Copy().Modify(mod.FlipX)

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
