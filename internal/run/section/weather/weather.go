package weather

import (
	"image/color"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
)

func NewWeatherSource() *particle.Source {
	g := particle.NewColorGenerator()
	return g.Generate(3)
}

func ClearWeather(ps *particle.Source) {
	ps.Generator = particle.NewColorGenerator(
		particle.Color(
			color.RGBA{50, 50, 50, 60},
			color.RGBA{10, 10, 10, 0},
			color.RGBA{50, 50, 50, 60},
			color.RGBA{10, 10, 10, 0},
		),
		particle.Shape(shape.Circle),
		particle.Size(intrange.NewLinear(4, 8)),
		particle.EndSize(intrange.NewLinear(4, 8)),
		particle.Spread(10, float64(oak.ScreenHeight)),
		particle.Speed(floatrange.NewLinear(.25, .75)),
		particle.Pos(float64(oak.ScreenWidth), 0),
		particle.Angle(floatrange.NewLinear(265, 275)),
		particle.Rotation(floatrange.NewLinear(-2, 2)),
		particle.LifeSpan(floatrange.NewLinear(300, 400)),
		particle.NewPerFrame(floatrange.NewLinear(1, 2)),
	)
}
