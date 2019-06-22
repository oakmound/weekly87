package vfx

import (
	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"

	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"image/color"
)

var (
	PushBack1 = func() particle.Generator {
		return particle.NewColorGenerator(
			particle.Color(color.RGBA{255, 158, 0, 255}, color.RGBA{0, 0, 0, 0},
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
)
