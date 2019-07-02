package enemies

import (
	"image"
	"path/filepath"

	"github.com/solovev/gopsd"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/render"
)

func initTree() {

	psdFilePath := filepath.Join("assets", "images", "64x64", "tree2.psd")
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
		// if c, ok := def.LayerColors[strings.ToLower(layer.Name)]; ok {

		// 	// Recolor this layer
		// 	sp.Filter(recolor.WithStrategy(recolor.ColorMix(c)))
		// }
		// Add this layer to the combined image
		// Todo: bug with shoulder having some pixel flashing
		combined.Append(sp)
	}
	combinedSp := combined.ToSprite()
	sh, err := render.MakeSheet(combinedSp.GetRGBA(), 62, 64, 0)
	dlog.ErrorCheck(err)
	sheet := sh.ToSprites()

	anims := map[string]render.Modifiable{}
	anims["standRT"] = sheet[0][0].Copy()
	anims["standLT"] = sheet[0][0].Copy().Modify(mod.FlipX)
	anims["walkRT"] = sheet[0][0].Copy()
	anims["walkLT"] = sheet[0][0].Copy().Modify(mod.FlipX)

	baseConstructor := Constructor{
		Dimensions:   floatgeom.Point2{32, 32},
		AnimationMap: anims,
		Speed:        floatgeom.Point2{3, 2},
		Bindings: map[string]func(*BasicEnemy, interface{}) int{
			"EnterFrame": func(b *BasicEnemy, frame interface{}) int {

				return 0
			},
		},
	}

	for size := 0; size < lastSize; size++ {
		for color := 0; color < lastColor; color++ {
			cons := baseConstructor.Copy()
			sizeVariants[size](cons)
			colorVariants[color](cons)
			setConstructor(int(Tree), size, color, cons)
		}
	}
}
