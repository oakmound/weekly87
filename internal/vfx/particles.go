package vfx

import (
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
)

var (
	PushBack1 = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(color.RGBA{50, 50, 170, 170}, color.RGBA{0, 0, 0, 0},
				color.RGBA{10, 10, 20, 20}, color.RGBA{0, 0, 0, 0}),
			particle.Shape(shape.Diamond),
			particle.Size(intrange.NewConstant(15)),
			particle.EndSize(intrange.NewConstant(4)),
			particle.Speed(floatrange.NewConstant(2)),
			particle.LifeSpan(floatrange.NewConstant(8)),
			particle.Spread(5, 10),
			particle.NewPerFrame(floatrange.NewConstant(4)),
			particle.Angle(floatrange.NewLinear(160, 200)),
		)
	}
	Death1 = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(color.RGBA{100, 158, 200, 255}, color.RGBA{0, 0, 0, 0},
				color.RGBA{125, 125, 125, 125}, color.RGBA{0, 0, 0, 0}),
			particle.Shape(shape.Diamond),
			particle.Size(intrange.NewConstant(10)),
			particle.EndSize(intrange.NewConstant(5)),
			particle.Speed(floatrange.NewConstant(2)),
			particle.LifeSpan(floatrange.NewConstant(2)),
			particle.Spread(5, 5),
			particle.NewPerFrame(floatrange.NewConstant(40)),
		)
	}
	WhiteSlash = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(
				color.RGBA{255, 255, 255, 255}, color.RGBA{},
				color.RGBA{255, 255, 255, 255}, color.RGBA{},
			),
			particle.Shape(shape.Square),
			particle.Size(intrange.NewLinear(3, 7)),
			particle.LifeSpan(floatrange.NewConstant(4)),
			particle.Speed(floatrange.NewLinear(2, 8)),
			particle.Angle(floatrange.NewLinear(215, 200)),
			particle.NewPerFrame(floatrange.NewConstant(3)),
		)
	}

	WhiteRing = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(
				color.RGBA{255, 255, 255, 255}, color.RGBA{},
				color.RGBA{255, 255, 255, 255}, color.RGBA{},
			),
			particle.Shape(shape.Square),
			particle.Size(intrange.NewLinear(3, 7)),
			particle.LifeSpan(floatrange.NewConstant(20)),
			particle.Speed(floatrange.NewConstant(4)),
			particle.Angle(floatrange.NewLinear(0, 360)),
			particle.NewPerFrame(floatrange.NewConstant(50)),
		)
	}

	RedRing = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(
				color.RGBA{255, 20, 20, 255}, color.RGBA{},
				color.RGBA{255, 20, 20, 255}, color.RGBA{},
			),
			particle.Shape(shape.Square),
			particle.Size(intrange.NewLinear(2, 6)),
			particle.LifeSpan(floatrange.NewConstant(20)),
			particle.Speed(floatrange.NewConstant(4)),
			particle.Angle(floatrange.NewLinear(0, 360)),
			particle.NewPerFrame(floatrange.NewConstant(10)),
		)
	}
)
