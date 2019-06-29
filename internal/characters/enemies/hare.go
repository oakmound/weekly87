package enemies

import (
	"math/rand"
	"path/filepath"

	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/oak/alg/floatgeom"
)

func initHare() {
	sheet, err := render.LoadSprites(filepath.Join("assets", "images"),
		filepath.Join("32x32", "Hare.png"), 32, 32, 0)
	dlog.ErrorCheck(err)
	anims := map[string]render.Modifiable{}
	anims["standRT"] = sheet[0][0].Copy()
	anims["standLT"] = sheet[0][0].Copy().Modify(mod.FlipX)
	anims["walkRT"] = render.NewSequence(4, sheet[0][0].Copy(), sheet[1][0].Copy())
	anims["walkLT"] = anims["walkRT"].Copy().Modify(mod.FlipX)

	baseConstructor := Constructor{
		Dimensions:   floatgeom.Point2{32, 32},
		AnimationMap: anims,
		Speed:        floatgeom.Point2{3, 2},
		Bindings: map[string]func(*BasicEnemy, interface{}) int{
			"EnterFrame": func(b *BasicEnemy, frame interface{}) int {
				f := frame.(int)
				// Simulate hops
				if f%52 == 0 {
					b.Speed = physics.NewVector(0, 0)
				} else if f%70 == 0 {
					b.Speed = physics.NewVector(
						-(rand.Float64()*b.baseSpeed.X()+1)*3,
						rand.Float64()*b.baseSpeed.Y()*float64(rand.Intn(2)*2-1),
					)
					if b.facing == "RT" {
						b.Speed.Scale(-1)
					}
				}
				return 0
			},
		},
	}

	// Todo: if this gets more complicated, we can have initHare return the baseConstructor
	// and do this in one place for all enemies
	for size := 0; size < lastSize; size++ {
		for color := 0; color < lastColor; color++ {
			cons := baseConstructor.Copy()
			sizeVariants[size](cons)
			colorVariants[color](cons)
			setConstructor(int(Hare), size, color, cons)
		}
	}
}
