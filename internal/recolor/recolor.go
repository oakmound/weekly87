package recolor

import (
	"image"
	"image/color"
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
