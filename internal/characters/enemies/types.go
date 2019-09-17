package enemies

import (
	"image/color"

	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/recolor"
)

type enemyType int

const (
	Hare   enemyType = iota
	Mantis enemyType = iota
	Tree   enemyType = iota
	TypeLimit
)

type Variant func(*Constructor)

const (
	baseColor = iota
	blueColor
	redColor
	blackColor
	purpleColor
	lastColor
)

const (
	baseSize = iota
	largeSize
	smallSize
	giantSize
	lastSize
)

var (
	colorVariants = [lastColor]Variant{
		baseColor: func(*Constructor) {},
		blueColor: func(c *Constructor) {
			for _, md := range c.AnimationMap {
				md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{100, 100, 255, 150})))
			}
			c.Speed = c.Speed.MulConst(.5)
		},
		redColor: func(c *Constructor) {
			for _, md := range c.AnimationMap {
				md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{255, 100, 100, 150})))
			}
			c.Speed = c.Speed.MulConst(1.5)
		},
		blackColor: func(c *Constructor) {
			for _, md := range c.AnimationMap {
				md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{50, 50, 50, 150})))
			}
			c.Health = 1000
		},
		purpleColor: func(c *Constructor) {
			for _, md := range c.AnimationMap {
				md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{255, 255, 100, 150})))
			}
			c.Health = 3
		},
	}
	sizeVariants = [lastSize]Variant{
		baseSize:  func(*Constructor) {},
		largeSize: changeSize(1.5),
		smallSize: changeSize(.5),
		giantSize: changeSize(2),
	}
)

func changeSize(mult float64) func(c *Constructor) {
	return func(c *Constructor) {
		c.Dimensions = c.Dimensions.MulConst(mult)
		if c.SpaceOffset != (physics.Vector{}) {
			c.SpaceOffset = c.SpaceOffset.Copy().Scale(mult)
		}
		for k, md := range c.AnimationMap {
			c.AnimationMap[k] = md.Modify(mod.Scale(mult, mult))

		}
	}
}

const VariantCount = lastSize * lastColor

var enemyTypeList = [TypeLimit]enemyType{
	Hare,
	Mantis,
	Tree,
}

func Init() {
	initHare()
	initMantis()
	initTree()
}
