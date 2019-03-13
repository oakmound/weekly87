package enemies

import (
	"math/rand"
	"path/filepath"

	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/render"

	"github.com/oakmound/oak/alg/floatgeom"
)

func initMantis() {
	sheet, err := render.LoadSprites(filepath.Join("assets", "images"),
		filepath.Join("32x32", "mantis.png"), 32, 32, 0)
	dlog.ErrorCheck(err)
	anims := map[string]render.Modifiable{}
	anims["standRT"] = sheet[0][0].Copy()
	anims["standLT"] = sheet[0][0].Copy().Modify(mod.FlipX)
	anims["walkRT"] = render.NewSequence(4, sheet[0][0].Copy(), sheet[1][0].Copy())
	anims["walkLT"] = anims["walkRT"].Copy().Modify(mod.FlipX)

	Constructors[Mantis] = Constructor{
		Dimensions:   floatgeom.Point2{32, 32},
		AnimationMap: anims,
		Speed: floatgeom.Point2{
			-1 * ((rand.Float64() * 4) + 1),
			-1 * ((rand.Float64() * 4) + 1),
		},
	}
}
