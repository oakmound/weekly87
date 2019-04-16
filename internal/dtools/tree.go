package dtools

import (
	"image/color"
	"image/draw"

	"github.com/oakmound/oak"

	"github.com/oakmound/oak/render"

	"github.com/oakmound/oak/collision"
)

// NewRTree creates a wrapper around a tree that supports coloring the spaces
func NewRTree(t *collision.Tree) *Rtree {
	rt := new(Rtree)
	rt.Tree = t
	rt.LayeredPoint = render.NewLayeredPoint(0, 0, -1)
	rt.OutlineColor = color.RGBA{0, 0, 255, 255}
	rt.ColorMap = map[collision.Label]color.RGBA{}
	return rt
}

type Rtree struct {
	*collision.Tree
	render.LayeredPoint
	OutlineColor color.RGBA
	ColorMap     map[collision.Label]color.RGBA
}

func (r *Rtree) GetDims() (int, int) {
	return oak.ScreenWidth, oak.ScreenHeight
}

func (r *Rtree) Draw(buff draw.Image) {
	r.DrawOffset(buff, 0, 0)
}

func (r *Rtree) DrawOffset(buff draw.Image, xOff, yOff float64) {
	// Get all spaces on screen
	screen := collision.NewUnassignedSpace(
		float64(oak.ViewPos.X),
		float64(oak.ViewPos.Y),
		float64(oak.ScreenWidth),
		float64(oak.ScreenHeight))
	hits := r.Tree.Hits(screen)
	// Draw spaces that are on screen (as outlines)
	for _, h := range hits {
		c := r.OutlineColor
		if found, ok := r.ColorMap[h.Label]; ok {
			c = found
		}
		for x := 0; x < int(h.GetW()); x++ {
			buff.Set(x+int(h.X()+xOff), int(h.Y()+yOff), c)
			buff.Set(x+int(h.X()+xOff), int(h.Y()+yOff)+int(h.GetH()), c)
		}
		for y := 0; y < int(h.GetH()); y++ {
			buff.Set(int(h.X()+xOff), y+int(h.Y()+yOff), c)
			buff.Set(int(h.X()+xOff)+int(h.GetW()), y+int(h.Y()+yOff), c)
		}
	}
}
