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
	"github.com/oakmound/oak/physics"
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

	SwordSwipe = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 90, 0, 255}),
		time.Second*4,
		func(u User) []characters.Character {
			dlog.Info("Trying to swipe at enemies")

			yDelta := 10.0
			xOffset := 100.0
			xDelta := 150.0
			if u.Direction() == "LT" {
				xOffset *= -1
				xDelta *= -1
			}

			pos := u.Vec().Copy()
			pos.Add(physics.NewVector(xOffset, 16.0))
			start := floatgeom.Point2{pos.X(), pos.Y() - yDelta}

			// Swipe
			// Should this be bigger?
			seq, err := render.LoadSheetSequence(filepath.Join("32x32", "BaseSlash.png"), 32, 32, 0, 32,
				0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1, 3, 1)
			dlog.ErrorCheck(err)

			chrs, err := Produce(
				StartAt(start),
				LineTo(start),
				FrameLength(16),
				FollowSpeed(u.GetDelta().Xp(), u.GetDelta().Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(seq),
			)

			dlog.ErrorCheck(err)
			return chrs
		},
	)

	HammerSmack = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 80, 80, 255}),
		time.Second*4,
		func(u User) []characters.Character {
			dlog.Info("Trying to swipe at enemies")

			yDelta := 10.0
			xOffset := 100.0
			xDelta := 150.0
			if u.Direction() == "LT" {
				xOffset *= -1
				xDelta *= -1
			}

			pos := u.Vec().Copy()
			pos.Add(physics.NewVector(xOffset, 16.0))
			start := floatgeom.Point2{pos.X(), pos.Y() - yDelta}

			// Swipe
			// Should this be bigger?
			seq, err := render.LoadSheetSequence(filepath.Join("32x32", "BaseSlash.png"), 32, 32, 0, 32,
				0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1, 3, 1)
			dlog.ErrorCheck(err)

			chrs, err := Produce(
				StartAt(start),
				LineTo(start),
				FrameLength(16),
				FollowSpeed(u.GetDelta().Xp(), u.GetDelta().Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(seq),
			)

			dlog.ErrorCheck(err)
			return chrs
		},
	)

	Rage = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{230, 5, 0, 255}),
		time.Second*5,
		func(u User) []characters.Character {
			fmt.Println("Just tried to rage a guy ", u)
			return nil
		},
	)

	SpearThrow = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 50, 150, 255}),
		time.Second*5,
		func(u User) []characters.Character {
			fmt.Println("Just tried to spear throw a guy ", u)
			return nil
		},
	)

	PartyShield = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{40, 200, 90, 255}),
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

			endDelta := 600.0
			if u.Direction() == "LT" {
				endDelta *= -1
			}
			end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
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

	SelfShield = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{110, 200, 110, 255}),
		time.Second*5,
		func(u User) []characters.Character {
			fmt.Println("To Be Implemented ", u)
			return nil
		},
	)
)
