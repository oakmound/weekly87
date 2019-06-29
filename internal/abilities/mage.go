package abilities

import (
	"image/color"
	"path/filepath"
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/abilities/buff"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/labels"
)

var (
	// FrostBolt is a simple projectile with slowing
	FrostBolt = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{10, 10, 200, 255}),
		time.Second*20,
		func(u User) []characters.Character {
			dlog.Info("Firing a frostbolt")
			pos := u.Vec()

			// Spell Display
			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{10, 10, 255, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(20)),
				particle.EndSize(intrange.NewConstant(3)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(30)),
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
				WithLabel(labels.EffectsEnemy),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

	//Fireball tries to cast a magical fire ball in front of the mage
	Fireball = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 200}),
		time.Second*10,
		func(u User) []characters.Character {
			dlog.Info("Firing a fireball")
			pos := u.Vec()

			// Spell Display
			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{255, 10, 10, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(5)),
				particle.EndSize(intrange.NewConstant(3)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(10)),
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
				WithLabel(labels.EffectsEnemy),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

	Blizzard = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{10, 10, 250, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			dlog.Info("making a blizzard")
			pos := u.Vec()

			// Spell Display
			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{10, 10, 255, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(10)),
				particle.EndSize(intrange.NewConstant(3)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(20)),
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
				WithLabel(labels.EffectsEnemy),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
	)

	FireWall = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			dlog.Info("Firing a fireball")
			pos := u.Vec()

			// Spell Display
			pg := particle.NewColorGenerator(
				particle.Color(color.RGBA{255, 10, 10, 255}, color.RGBA{0, 0, 0, 0},
					color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
				particle.Shape(shape.Diamond),
				particle.Size(intrange.NewConstant(10)),
				particle.EndSize(intrange.NewConstant(3)),
				particle.Speed(floatrange.NewConstant(1)),
				particle.LifeSpan(floatrange.NewConstant(20)),
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
				WithLabel(labels.EffectsEnemy),
			)
			dlog.ErrorCheck(err)
			return chrs
		},
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
				//ArcTo(end),
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
				//ArcTo(end),
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
