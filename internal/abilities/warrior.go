package abilities

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
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

func thwack(image string, xOffset, yDelta float64, hitEffects map[string]float64, mods ...mod.Mod) func(User) []characters.Character {
	hits := map[collision.Label]collision.OnHit{
		labels.Enemy: func(a, b *collision.Space) {
			b.CID.Trigger("Attacked", hitEffects)
		},
	}
	var md render.Modifiable
	seq, err := render.LoadSheetSequence(image, 32, 32, 0, 32,
		0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1, 3, 1)

	dlog.ErrorCheck(err)
	return func(u User) []characters.Character {
		md = seq.Copy()
		dlog.Info("Trying to swipe at enemies")
		if u.Direction() == "LT" {
			xOffset *= -1
			md = seq.Copy().Modify(mod.FlipX)
		}

		pos := u.Vec().Copy()
		pos.Add(physics.NewVector(xOffset, 16.0))
		start := floatgeom.Point2{pos.X(), pos.Y() - yDelta}

		for _, m := range mods {
			md = md.Modify(m)
		}
		chrs, err := Produce(
			StartAt(start),
			LineTo(start),
			FrameLength(16),
			FollowSpeed(u.GetDelta().Xp(), u.GetDelta().Yp()),
			WithHitEffects(hits),
			WithLabel(labels.EffectsEnemy),
			WithRenderable(md),
		)
		dlog.ErrorCheck(err)
		sfx.Play("slashHeavy")
		return chrs
	}
}

// Warrior abilities
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

	// SwordSwipe is a basic thwack
	SwordSwipe = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{150, 150, 0, 200}), slashIcon),
		time.Second*4,
		thwack(filepath.Join("32x32", "BaseSlash.png"), 100, 10, dmg),
	)

	// HammerSmack is currently a huge thwack TODO: have new animation with hammer
	HammerSmack = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{150, 80, 120, 255}), hammerIcon),

		time.Second*8,
		thwack(filepath.Join("32x32", "BaseSlash.png"), 100, 26, map[string]float64{"damage": 1.0, "pushback": 120.0}, mod.Scale(2, 2)),
	)

	// Rage is a multistrike attack that impacts party movement
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
				WithHitEffects(baseHit),
				WithLabel(labels.EffectsEnemy),
				WithRenderable(down.Copy()),
				PlaySFX("slashHeavy"),
			)(Producer{})

			hit3 := And(
				StartAt(floatgeom.Point2{xOffset / 2, 0}),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithHitEffects(baseHit),
				WithRenderable(up.Copy()),
				WithLabel(labels.EffectsEnemy),
				Then(Chain(hit4)),
				PlaySFX("slashLight"),
			)(Producer{})

			hit2 := And(
				StartAt(floatgeom.Point2{0, yDelta * 5}),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithHitEffects(baseHit),
				WithRenderable(up.Copy()),
				WithLabel(labels.EffectsEnemy),
				Then(Chain(hit3)),
				PlaySFX("slashHeavy"),
			)(Producer{})

			chrs, err := Produce(
				StartAt(start),
				LineTo(start),
				FrameLength(16),
				FollowSpeed(delta.Xp(), delta.Yp()),
				WithHitEffects(baseHit),
				WithRenderable(down),
				WithLabel(labels.EffectsEnemy),
				Then(Chain(hit2)),
				PlaySFX("slashLight"),
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

	// Party Shield is a slower moving buff that protects the whole party
	PartyShield = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{40, 200, 90, 255}), shieldAuraIcon),
		time.Second*15,
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

			endDelta := 680.0
			if u.Direction() == "LT" {
				endDelta *= -1
			}
			end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
			chrs, err := Produce(
				StartAt(floatgeom.Point2{pos.X(), pos.Y()}),
				//ArcTo(end),
				LineTo(end),
				WithParticles(pg),
				Then(AndDo(Drop(banner), DoPlay("bannerPlaced1"))),
				FollowSpeed(u.GetDelta().Xp(), nil),
				PlaySFX("warriorCast1"),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

	// SelfShield is a fast flying single person shield
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
				particle.Color(color.RGBA{120, 255, 0, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(10)),
				particle.EndSize(intrange.NewConstant(5)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(15)),
			)

			endDelta := 520.0
			if u.Direction() == "LT" {
				endDelta *= -1
			}
			chrs, err := Produce(
				StartAt(floatgeom.Point2{pos.X(), pos.Y()}),

				LineTo(floatgeom.Point2{pos.X() + endDelta, pos.Y()}),
				FollowSpeed(u.GetDelta().Xp(), nil),
				WithParticles(pg),
				Then(AndDo(Drop(banner), DoPlay("bannerPlaced1"))),
				FrameLength(30),
				PlaySFX("warriorCast1"),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)
}
