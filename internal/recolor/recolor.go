package recolor

import (
	"image"
	"image/color"
	"time"

	"math/rand"

	"sort"
)

func Recolor(colors map[color.RGBA]color.RGBA) func(*image.RGBA) {
	return func(rgba *image.RGBA) {
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				if c2, ok := colors[rgba.RGBAAt(x, y)]; ok {
					rgba.SetRGBA(x, y, c2)
				}
			}
		}
	}
}

// WithStrategy finds all the colors in an image, then applies the strategy given
// to provide a recolor function.
func WithStrategy(strat func([]color.RGBA) []color.RGBA) func(*image.RGBA) {
	return func(rgba *image.RGBA) {
		bounds := rgba.Bounds()
		w := bounds.Max.X
		h := bounds.Max.Y
		colorM := make(map[color.RGBA]struct{})
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				colorM[rgba.RGBAAt(x, y)] = struct{}{}
			}
		}
		colors := make([]color.RGBA, 0, len(colorM))
		for c := range colorM {
			colors = append(colors, c)
		}

		newColors := strat(colors)

		colorMap := make(map[color.RGBA]color.RGBA, len(newColors))
		for i, v := range colors {
			colorMap[v] = newColors[i]
		}
		Recolor(colorMap)(rgba)
	}
}

// StrategyA iterates through all individual r, g, b values and reassigns them
// randomly. This mantains that different colors with similar components will
// keep being similar to one another.
// Values will be rounded to one another within the provided rounding bridge.
// Todo: name
// If seed is provided, will seed math/rand at start
func StrategyA(bridge uint8, seed bool) func([]color.RGBA) []color.RGBA {
	return func(colors []color.RGBA) []color.RGBA {
		if seed {
			rand.Seed(time.Now().UnixNano())
		}
		reds := make([]uint8, 0)
		greens := make([]uint8, 0)
		blues := make([]uint8, 0)
		for _, c := range colors {
			reds = append(reds, c.R)
			greens = append(reds, c.G)
			blues = append(reds, c.B)
		}

		sort.Slice(reds, func(i, j int) bool { return reds[i] < reds[j] })
		sort.Slice(greens, func(i, j int) bool { return greens[i] < greens[j] })
		sort.Slice(blues, func(i, j int) bool { return blues[i] < blues[j] })

		mReds := make(map[uint8]uint8)
		mBlues := make(map[uint8]uint8)
		mGreens := make(map[uint8]uint8)

		last := uint8(0)
		for _, r := range reds {
			if last != 0 && last+bridge > r {
				mReds[r] = last
			} else {
				mReds[r] = r
				last = r
			}
		}
		last = uint8(0)
		for _, g := range greens {
			if last != 0 && last+bridge > g {
				mGreens[g] = last
			} else {
				mGreens[g] = g
				last = g
			}
		}
		last = uint8(0)
		for _, b := range blues {
			if last != 0 && last+bridge > b {
				mBlues[b] = last
			} else {
				mBlues[b] = b
				last = b
			}
		}

		convRed := make(map[uint8]uint8)
		convGreen := make(map[uint8]uint8)
		convBlue := make(map[uint8]uint8)

		for _, v := range mReds {
			convRed[v] = uint8(rand.Float64() * 255)
		}
		for _, v := range mGreens {
			convGreen[v] = uint8(rand.Float64() * 255)
		}
		for _, v := range mBlues {
			convBlue[v] = uint8(rand.Float64() * 255)
		}

		outColors := make([]color.RGBA, len(colors))
		for i, c := range colors {
			outColors[i] = color.RGBA{
				convRed[mReds[c.R]],
				convGreen[mReds[c.G]],
				convBlue[mReds[c.B]],
				c.A,
			}
		}
		return outColors
	}
}

// ColorShift moves all colors towards the given inColor's R G and B
// by taking inColor.A and combining inColor proportionally with a brightness-scaled
// version of the color to be converted
func ColorShift(inColor color.RGBA) func([]color.RGBA) []color.RGBA {
	return func(colors []color.RGBA) []color.RGBA {

		inColorFactor := float64(inColor.A) / 255

		inBrightness := float64(inColor.R) + float64(inColor.G) + float64(inColor.B)
		inRBightness := float64(inColor.R) / inBrightness
		inGBrightness := float64(inColor.G) / inBrightness
		inBBrightness := float64(inColor.B) / inBrightness

		inColorComponent := [3]float64{
			float64(inColor.R) * inColorFactor,
			float64(inColor.G) * inColorFactor,
			float64(inColor.B) * inColorFactor,
		}

		outColors := make([]color.RGBA, len(colors))

		for i, v := range colors {
			if v.A == 0 {
				continue
			}
			brightness := float64(v.R) + float64(v.G) + float64(v.B)
			// Trying not to recolor pure black (outlines)
			if brightness == 0 {
				continue
			}
			rBrightness := (brightness * inRBightness * (1 - inColorFactor)) + inColorComponent[0]
			gBrightness := (brightness * inGBrightness * (1 - inColorFactor)) + inColorComponent[1]
			bBrightness := (brightness * inBBrightness * (1 - inColorFactor)) + inColorComponent[2]

			outColors[i] = color.RGBA{uint8(rBrightness), uint8(gBrightness), uint8(bBrightness), v.A}
		}

		return outColors
	}
}

// ColorMix moves all colors towards the given inColor's R G and B
//
func ColorMix(inColor color.RGBA) func([]color.RGBA) []color.RGBA {
	return func(colors []color.RGBA) []color.RGBA {

		inColorFactor := float64(inColor.A) / 255

		// inBrightness := float64(inColor.R) + float64(inColor.G) + float64(inColor.B)
		// inRBightness := float64(inColor.R) / inBrightness
		// inGBrightness := float64(inColor.G) / inBrightness
		// inBBrightness := float64(inColor.B) / inBrightness

		inColorComponent := [3]float64{
			float64(inColor.R) * inColorFactor,
			float64(inColor.G) * inColorFactor,
			float64(inColor.B) * inColorFactor,
		}

		outColors := make([]color.RGBA, len(colors))

		for i, v := range colors {
			if v.A == 0 {
				continue
			}
			brightness := float64(v.R) + float64(v.G) + float64(v.B)
			// Trying not to recolor pure black (outlines)
			if brightness == 0 {
				continue
			}

			// Goals
			// take into account what we are shifting towards
			// take into account alpha of origin but also shift (example white mage?)
			// take into account somewhat the basic colors of the source

			rBrightness := (float64(v.R) * (1 - inColorFactor)) + inColorComponent[0]
			gBrightness := (float64(v.G) * (1 - inColorFactor)) + inColorComponent[1]
			bBrightness := (float64(v.B) * (1 - inColorFactor)) + inColorComponent[2]

			// rBrightness := ((float64(v.R)/brightness + inRBightness*2) / 3 * (1 - inColorFactor)) + inColorComponent[0]
			// gBrightness := ((float64(v.G)/brightness + inGBrightness*2) / 3 * (1 - inColorFactor)) + inColorComponent[1]
			// bBrightness := ((float64(v.B)/brightness + inBBrightness*2) / 3 * (1 - inColorFactor)) + inColorComponent[2]

			outColors[i] = color.RGBA{uint8(rBrightness), uint8(gBrightness), uint8(bBrightness), v.A}
		}

		return outColors
	}
}
