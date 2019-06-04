package abilities

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/labels"
)

var (
	//SpearStab tries to stab in front of the warrior
	SpearStab = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 200, 0, 255}),
		time.Second*5,
		func(u User) []characters.Character {
			fmt.Println("Just tried to stab a guy ", u)
			return nil
		},
	)

	Shield = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{40, 200, 125, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			pos := u.Vec()

			animFilePath := (filepath.Join("16x32", "banner.png"))
			seq, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 5, []int{0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1}...)
			dlog.ErrorCheck(err)

			if err != nil {
				dlog.Error(err)
				return nil
			}

			buffIcon, err2 := render.LoadSprite(filepath.Join("assets/images", "16x16"), "place_holder_buff.png")
			dlog.ErrorCheck(err2)

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Shield(buffIcon, 20*time.Second, 2)))(Producer{})

			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{255, 255, 0, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(10)),
				particle.EndSize(intrange.NewConstant(5)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(15)),
			)

			end := floatgeom.Point2{pos.X() + 600, pos.Y()}
			chrs, err := Produce(
				StartAt(floatgeom.Point2{pos.X(), pos.Y()}),
				//ArcTo(end),
				LineTo(end),
				WithParticles(pg),
				Then(Drop(banner)),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)
)
