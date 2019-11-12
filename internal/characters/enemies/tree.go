package enemies

import (
	"image"
	"image/color"
	"path/filepath"
	"strings"
	"io/ioutil"

	"github.com/solovev/gopsd"

	"github.com/oakmound/weekly87/internal/recolor"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/fileutil"
)

func initTree() {

	psdFilePath := filepath.Join("assets", "images", "64x64", "tree2.psd")
	rd, err := fileutil.Open(psdFilePath)
	dlog.ErrorCheck(err)
	data, err := ioutil.ReadAll(rd)
	dlog.ErrorCheck(err)
	psd, err := gopsd.ParseFromBuffer(data)
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
	combined.Append(render.NewEmptySprite(0, 0, int(psd.Width), int(psd.Height))) // Make sure this is here in case there is no layer that encompasses the whole thing
	combinedSp := combined.Slice(combined.Len()-2, combined.Len()).ToSprite()
	sh, err := render.MakeSheet(combinedSp.GetRGBA(), 64, 64, 0)
	dlog.ErrorCheck(err)
	sheet := sh.ToSprites()

	anims := map[string]render.Modifiable{}
	anims["standRT"] = sheet[0][0].Copy()
	anims["standLT"] = sheet[0][0].Copy().Modify(mod.FlipX)
	anims["walkRT"] = sheet[0][0].Copy()
	anims["walkLT"] = sheet[0][0].Copy().Modify(mod.FlipX)

	baseConstructor := Constructor{
		Dimensions:   floatgeom.Point2{20, 50},
		AnimationMap: anims,
		Speed:        floatgeom.Point2{0, 0},
		SpaceOffset:  physics.NewVector(-22, -6),
		Bindings: map[string]func(*BasicEnemy, interface{}) int{
			"EnterFrame": func(b *BasicEnemy, frame interface{}) int {

				return 0 
			}, 
		},   
		Health: 2,
	}

	toOverwriteWith := combined.Slice(0,combined.Len()-2)
	toOverwriteWithFlipped := toOverwriteWith.Copy().Modify(mod.FlipX).(*render.CompositeM)
	for size := 0; size < lastSize; size++ {
		for col := 0; col < lastColor; col++ {
			cons := baseConstructor.Copy()
			colorVariants[col](cons)
			if size == 0 && col == 0 {
				for _, md := range cons.AnimationMap {
					md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{140, 200, 140, 100})))
				}
			}
			// Overwrite with non-zero layers from combined
			for k, md := range cons.AnimationMap {
				var cmp *render.CompositeM = toOverwriteWith
				if strings.HasSuffix(k, "LT") {
					cmp = toOverwriteWithFlipped
				}
				cmp.Append(md)
				sprite := cmp.ToSprite()
				cons.AnimationMap[k] = sprite
			}
			sizeVariants[size](cons)
			setConstructor(int(Tree), size, col, cons)
		}
	}
}
