package abilities

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

var (
	//Fireball tries to cast a magical fire ball in front of the mage
	Fireball = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			fmt.Println("Just tried to burn a guy ", u)
			return nil
		},
	)
	Invulnerability = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{255, 255, 125, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			pos := u.Vec()

			sp, err := render.LoadSprite("images", filepath.Join("16x32", "banner.png"))
			if err != nil {
				dlog.Error(err)
				return nil
			}

			banner := WithRenderable(sp)(Producer{})

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
