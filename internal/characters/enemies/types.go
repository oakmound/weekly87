package enemies

import (
	"image/color"

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
				md.Filter(recolor.WithStrategy(recolor.ColorMix(color.RGBA{100, 255, 100, 150})))
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
		baseSize: func(*Constructor) {},
		largeSize: func(c *Constructor) {
			c.Dimensions = c.Dimensions.MulConst(1.5)
			for k, md := range c.AnimationMap {
				c.AnimationMap[k] = md.Modify(mod.Scale(1.5, 1.5))
			}
		},
		smallSize: func(c *Constructor) {
			c.Dimensions = c.Dimensions.MulConst(.5)
			for k, md := range c.AnimationMap {
				c.AnimationMap[k] = md.Modify(mod.Scale(.5, .5))
			}
		},
		giantSize: func(c *Constructor) {
			c.Dimensions = c.Dimensions.MulConst(2)
			for k, md := range c.AnimationMap {
				c.AnimationMap[k] = md.Modify(mod.Scale(2, 2))
			}
		},
	}
)

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
