package abilities

import (
	"image/color"
	"path/filepath"
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/abilities/buff"

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

func storm(speed floatrange.Range, sc, ec color.Color, xSpreadFactor float64, opts func(particle.Generator),
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
			WithLabel(labels.EffectsEnemy),
			FollowSpeed(delta.Xp(), nil),
		)
		dlog.ErrorCheck(err)
		return chrs
	}
}

var (
	baseHit = map[collision.Label]collision.OnHit{
		labels.Enemy: func(a, b *collision.Space) {
			b.CID.Trigger("Attacked", map[string]float64{"damage": 1.0})
		},
	}

	frostHit = map[collision.Label]collision.OnHit{
		labels.Enemy: func(a, b *collision.Space) {
			b.CID.Trigger("Attacked", map[string]float64{"frost": 5.0})
		},
	}

	// FrostBolt is a simple projectile with slowing
	FrostBolt = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{10, 10, 200, 255}),
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
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 200}),
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

	// Blizzard ice storm
	Blizzard = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{10, 10, 250, 255}),
		time.Second*10,
		storm(floatrange.NewLinear(3, 8), color.RGBA{10, 10, 255, 255}, color.RGBA{125, 125, 125, 125}, 1.5, particle.And(), map[string]float64{"frost": 1.2}),
	)

	// Firewall
	FireWall = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Second*1,
		storm(floatrange.NewLinear(3, 8), color.RGBA{255, 10, 10, 255}, color.RGBA{125, 125, 125, 125}, 1, particle.And(particle.NewPerFrame(floatrange.NewLinear(2, 4)), particle.Size(intrange.NewConstant(20))), map[string]float64{"damage": 1}),
	)

	// Rez
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
	// Invulnerability
	Invulnerability = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{90, 240, 90, 255}),
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

	// GameBreakerFireBall is debug murder ability
	GameBreakerFireBall = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Millisecond*1,
		func(u User) []characters.Character {
			dlog.Info("Firing a fireball")
			pos := u.Vec()

			// Spell Display
			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{255, 10, 10, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(50)),
				particle.EndSize(intrange.NewConstant(50)),
				particle.Speed(floatrange.NewConstant(3)),
				particle.LifeSpan(floatrange.NewConstant(1)),
			)
			// endDelta := 1200.0
			endDelta := float64(oak.ScreenWidth * 1)
			if u.Direction() == "LT" {
				endDelta *= -1
			}
			end := floatgeom.Point2{pos.X() + endDelta, pos.Y()}
			chrs, err := Produce(
				StartAt(floatgeom.Point2{pos.X(), pos.Y()}),
				FrameLength(400),
				LineTo(end),
				WithParticles(pg),
				WithLabel(labels.EffectsEnemy),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)
)
