package abilities

import (
	"image/color"
	"path/filepath"
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/sfx"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/labels"
)

func bolt(image string, frames int, endDelta float64, opts func(particle.Generator),
	hitEffects map[string]float64) func(u User) []characters.Character {
	return func(u User) []characters.Character {
		dlog.Info("Firing a bolt!")
		pos := u.Vec()

		// Spell Display
		r, err := render.LoadSheetSequence(image, 16, 16, 0, 16,
			0, 0, 1, 0, 0, 1, 1, 1)
		dlog.ErrorCheck(err)

		pg := particle.NewCollisionGenerator(
			particle.NewColorGenerator(opts),
			particle.Fragile(true),
			particle.HitMap(map[collision.Label]collision.OnHit{
				labels.Enemy: func(a, b *collision.Space) {
					b.CID.Trigger("Attacked", hitEffects)
				},
			}),
		)
		if u.Direction() == "LT" {
			endDelta *= -1
		}
		delta := u.GetDelta()
		end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
		chrs, err := Produce(
			StartAt(floatgeom.Point2{pos.X(), pos.Y()}),
			LineTo(end),
			FrameLength(frames),
			WithParticles(pg),
			WithRenderable(r),
			FollowSpeed(delta.Xp(), nil),
		)
		dlog.ErrorCheck(err)
		return chrs
	}
}

func storm(image string, dur time.Duration, speed floatrange.Range, sc, ec color.Color, xSpreadFactor float64, opts func(particle.Generator),
	hitEffects map[string]float64) func(u User) []characters.Character {
	// Spell Display

	return func(u User) []characters.Character {
		rs, err := render.LoadSprites("", image, 16, 16, 0)
		dlog.ErrorCheck(err)
		delta := u.GetDelta()
		pg := particle.NewSpriteGenerator(
			particle.Sprite(rs[0][0]),
			particle.Angle(floatrange.NewLinear(240, 300)),
			particle.Size(intrange.NewConstant(10)),
			particle.EndSize(intrange.NewConstant(3)),
			particle.Speed(speed),
			particle.NewPerFrame(floatrange.NewLinear(2, 7)),
			particle.LifeSpan(floatrange.NewLinear(200, 201)),
			particle.Spread(float64(oak.ScreenWidth)*xSpreadFactor, 0),
			opts,
		)
		endDelta := 600.0
		if u.Direction() == "LT" {
			endDelta *= -1
		}

		cpg := particle.NewCollisionGenerator(
			pg,
			particle.Fragile(true),
			particle.HitMap(map[collision.Label]collision.OnHit{
				labels.Enemy: func(a, b *collision.Space) {
					b.CID.Trigger("Attacked", hitEffects)
				},
			}),
		)

		// end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
		chrs, err := Produce(
			StartAt(floatgeom.Point2{float64(oak.ViewPos.X), 0}),
			WithParticles(cpg),

			FollowSpeed(delta.Xp(), nil),
			Duration(dur),
		)
		dlog.ErrorCheck(err)
		sfx.Play("stormEffect")
		return chrs
	}
}

func shower(speed floatrange.Range, dur time.Duration, sc, ec color.Color, xSpreadFactor float64, opts func(particle.Generator),
	hitEffects map[string]float64) func(u User) []characters.Character {
	return func(u User) []characters.Character {
		delta := u.GetDelta()
		pg := particle.NewColorGenerator(
			particle.Angle(floatrange.NewLinear(240, 300)),
			particle.Color(sc, color.RGBA{0, 0, 0, 0},
				ec, color.RGBA{0, 0, 0, 0}),
			particle.Shape(shape.Diamond),
			particle.Size(intrange.NewConstant(10)),
			particle.EndSize(intrange.NewConstant(3)),
			particle.Speed(speed),
			particle.NewPerFrame(floatrange.NewLinear(2, 7)),
			particle.LifeSpan(floatrange.NewLinear(200, 201)),
			particle.Spread(float64(oak.ScreenWidth)*xSpreadFactor, 0),
			opts,
		)
		endDelta := 600.0
		if u.Direction() == "LT" {
			endDelta *= -1
		}

		cpg := particle.NewCollisionGenerator(
			pg,
			particle.Fragile(true),
			particle.HitMap(map[collision.Label]collision.OnHit{
				labels.Enemy: func(a, b *collision.Space) {
					b.CID.Trigger("Attacked", hitEffects)
				},
			}),
		)

		// end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
		chrs, err := Produce(
			StartAt(floatgeom.Point2{float64(oak.ViewPos.X), 0}),
			WithParticles(cpg),
			FollowSpeed(delta.Xp(), nil),
			Duration(dur),
		)
		dlog.ErrorCheck(err)
		return chrs
	}
}

var (
	FrostBolt, Fireball, Blizzard, FireWall, FireStorm, Rez, Invulnerability, Slow, CooldownRework, GameBreakerFireBall *ability

	frostHit = map[collision.Label]collision.OnHit{
		labels.Enemy: func(a, b *collision.Space) {
			b.CID.Trigger("Attacked", map[string]float64{"frost": 5.0})
		},
	}
)

func MageInit() {

	// FrostBolt is a simple projectile with slowing
	FrostBolt = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{10, 10, 200, 200}), blueBlastIcon),

		time.Second*3,
		bolt(filepath.Join("16x16", "icebolt.png"),
			200,
			600,
			particle.And(
				particle.LifeSpan(floatrange.NewSpread(30, 5)),
				particle.Spread(4, 4),
				particle.Color(color.RGBA{150, 150, 255, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewLinear(20, 10)),
				particle.EndSize(intrange.NewConstant(3)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.Pos(8, 8),
			),
			map[string]float64{"frost": 5.0},
		),
	)

	//Fireball tries to cast a magical fire ball in front of the mage
	Fireball = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{200, 160, 160, 200}), redBlastIcon),
		time.Second*10,
		bolt(filepath.Join("16x16", "fireball.png"),
			200,
			600,
			particle.And(
				particle.NewPerFrame(floatrange.NewSpread(4, 2)),
				particle.LifeSpan(floatrange.NewSpread(30, 5)),
				particle.Speed(floatrange.NewSpread(4, .5)),
				particle.Spread(6, 6),
				particle.Color(
					color.RGBA{255, 155, 155, 255},
					color.RGBA{10, 50, 50, 0},
					color.RGBA{255, 100, 60, 255},
					color.RGBA{0, 10, 10, 0},
				),
				// particle.Color2(
				// 	color.RGBA{255, 230, 220, 255},
				// 	color.RGBA{10, 50, 50, 0},
				// 	color.RGBA{120, 50, 40, 140},
				// 	color.RGBA{20, 20, 20, 0},
				// ),
				particle.Size(intrange.NewSpread(10, 5)),
				particle.Shape(shape.Circle),
				//particle.Progress(render.CircularProgress),
				particle.Pos(8, 8),
			),
			map[string]float64{"damage": 1.0},
		),
	)

	// Blizzard creates a large slowing icestorm
	Blizzard = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{10, 10, 200, 200}), blueBlastDIcon),
		time.Second*10,
		shower(floatrange.NewLinear(3, 8), time.Second*3, color.RGBA{10, 10, 255, 255}, color.RGBA{125, 125, 125, 125}, 1.5, particle.And(), map[string]float64{"frost": 1.2}),
	)

	// FireWall is a short lived long cooldown vertical destructive force
	FireWall = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 200}), redBlastDIcon),
		time.Second*20,
		shower(floatrange.NewLinear(3, 8), time.Second*2, color.RGBA{255, 10, 10, 255}, color.RGBA{125, 125, 125, 125}, 1, particle.And(particle.NewPerFrame(floatrange.NewLinear(2, 4)), particle.Size(intrange.NewConstant(20))), map[string]float64{"damage": 1}),
	)

	// FireStorm is a short lived long cooldown vertical destructive force
	FireStorm = NewAbility(
		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 200}), redBlastDIcon),
		time.Second*20,
		storm(filepath.Join("16x16", "fireball.png"),
			time.Second*2,
			floatrange.NewLinear(3, 8), color.RGBA{255, 10, 10, 255}, color.RGBA{125, 125, 125, 125}, 2,
			particle.And(particle.NewPerFrame(floatrange.NewLinear(0, 2)), particle.Size(intrange.NewConstant(20))), map[string]float64{"damage": 1}),
	)

	// Rez the first person who is dead in the party on pickup
	Rez = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 200, 200, 255}),
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

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Rez))(Producer{})

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
				LineTo(end),
				WithParticles(pg),
				Then(Drop(banner)),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)
	// Invulnerability gives a temp buff of near invuln to the entire living party
	Invulnerability = NewAbility(

		render.NewCompositeM(render.NewColorBox(64, 64, color.RGBA{90, 240, 90, 255}), shieldIcon),
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

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Invulnerable(render.NewColorBox(BuffIconSize, BuffIconSize, color.RGBA{250, 250, 0, 255}), 6*time.Second)))(Producer{})

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
				Then(
					AndDo(Drop(banner), DoPlay("bannerPlaced1")),
				),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

	// Slow
	Slow = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{120, 120, 120, 255}),
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

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Invulnerable(render.NewColorBox(BuffIconSize, BuffIconSize, color.RGBA{250, 250, 0, 255}), 6*time.Second)))(Producer{})

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
	// CooldownRework
	CooldownRework = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{120, 120, 120, 255}),
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

			banner := And(WithRenderable(seq),
				WithLabel(labels.EffectsPlayer),
				WithBuff(buff.Invulnerable(render.NewColorBox(BuffIconSize, BuffIconSize, color.RGBA{250, 250, 0, 255}), 6*time.Second)))(Producer{})

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
				LineTo(end),
				WithParticles(pg),
				Then(Drop(banner)),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

}
