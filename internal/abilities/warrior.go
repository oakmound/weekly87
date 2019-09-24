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
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/sfx"
)

func thwack(image string, xOffset, yDelta float64, mods ...mod.Mod) func(User) []characters.Character {
	return func(u User) []characters.Character {

		var md render.Modifiable
		seq, err := render.LoadSheetSequence(image, 32, 32, 0, 32,
			0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1, 3, 1)
		md = seq

		dlog.Info("Trying to swipe at enemies")
		if u.Direction() == "LT" {
			xOffset *= -1
			md = seq.Copy().Modify(mod.FlipX)
		}

		pos := u.Vec().Copy()
		pos.Add(physics.NewVector(xOffset, 16.0))
		start := floatgeom.Point2{pos.X(), pos.Y() - yDelta}

		dlog.ErrorCheck(err)
		for _, m := range mods {
			md = md.Modify(m)
		}
		chrs, err := Produce(
			StartAt(start),
			LineTo(start),
			FrameLength(16),
			FollowSpeed(u.GetDelta().Xp(), u.GetDelta().Yp()),
			WithLabel(labels.EffectsEnemy),
			WithRenderable(md),
		)
		dlog.ErrorCheck(err)
		sfx.Play("slashHeavy")
		return chrs
	}
}

var (
	SpearStab, SwordSwipe, HammerSmack, Rage, SpearThrow, PartyShield, SelfShield *ability
)

// WarriorInit run by abilities to set up the ability attributes
func WarriorInit() {
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
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{150, 150, 0, 200}), slashIcon),
		time.Second*4,
		thwack(filepath.Join("32x32", "BaseSlash.png"), 100, 10),
	)

	HammerSmack = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 80, 80, 255}),
		time.Second*8,
		thwack(filepath.Join("32x32", "BaseSlash.png"), 100, 26, mod.Scale(2, 2)),
	)

	Rage = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{200, 5, 0, 200}), downSlashIcon),
		time.Second*5,
		func(u User) []characters.Character {
			var down render.Modifiable
			var err error
			// For efficiency in the future, we could pre-load these assets
			down, err = render.LoadSheetSequence(filepath.Join("32x32", "BaseSlash.png"), 32, 32, 0, 32,
				0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1, 3, 1)
			up := down.Copy().Modify(mod.FlipY)

			xOffset := float64(100)
			yDelta := float64(16)

			dlog.Info("Trying to swipe at enemies")
			if u.Direction() == "LT" {
				xOffset *= -1
				down = down.Copy().Modify(mod.FlipX)
				up = up.Copy().Modify(mod.FlipX)
			}

			pos := u.Vec().Copy()
			pos.Add(physics.NewVector(xOffset, 16.0))
			start := floatgeom.Point2{pos.X(), pos.Y() - yDelta}

			dlog.ErrorCheck(err)

			delta := u.GetDelta()

			hit4 := And(
				StartAt(floatgeom.Point2{0, -yDelta * 6}),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(down.Copy()),
			)(Producer{})

			hit3 := And(
				StartAt(floatgeom.Point2{xOffset / 2, 0}),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(up.Copy()),
				Then(Chain(hit4)),
			)(Producer{})

			hit2 := And(
				StartAt(floatgeom.Point2{0, yDelta * 5}),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(up.Copy()),
				Then(Chain(hit3)),
			)(Producer{})

			chrs, err := Produce(
				StartAt(start),
				LineTo(start),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(down),
				Then(Chain(hit2)),
			)
			dlog.ErrorCheck(err)
			event.Trigger("RageStart", nil)
			return chrs
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
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{40, 200, 90, 255}), shieldAuraIcon),
		time.Second*10,
		func(u User) []characters.Character {
			pos := u.Vec()

			animFilePath := (filepath.Join("16x32", "banner.png"))
			seq, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 5, []int{0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1}...)
			dlog.ErrorCheck(err)

			buffIcon, err2 := render.LoadSprite(filepath.Join("assets/images", "16x16"), "place_holder_buff.png")
			dlog.ErrorCheck(err2)

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Shield(buffIcon, 20*time.Second, 2, false)))(Producer{})

			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{0, 255, 255, 255}, color.RGBA{0, 0, 0, 0},
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
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{110, 200, 110, 255}), shieldIcon),
		time.Second*5,
		func(u User) []characters.Character {
			pos := u.Vec()

			animFilePath := (filepath.Join("16x32", "banner.png"))
			seq, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 5, []int{0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1}...)
			dlog.ErrorCheck(err)

			buffIcon, err2 := render.LoadSprite(filepath.Join("assets/images", "16x16"), "place_holder_buff.png")
			dlog.ErrorCheck(err2)

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Shield(buffIcon, 20*time.Second, 5, true)))(Producer{})

			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{255, 255, 0, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(10)),
				particle.EndSize(intrange.NewConstant(5)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(15)),
			)

			endDelta := 400.0
			if u.Direction() == "LT" {
				endDelta *= -1
			}
			chrs, err := Produce(
				StartAt(floatgeom.Point2{pos.X(), pos.Y()}),
				//ArcTo(end),
				LineTo(floatgeom.Point2{pos.X() + endDelta, pos.Y()}),
				WithParticles(pg),
				Then(Drop(banner)),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)
}
