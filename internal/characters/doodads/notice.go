package doodads

import (
	"image/color"
	"math"
	"math/rand"
	"strings"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/shiny/materialdesign/colornames"
)

// NewNote creates overlapping notes with random colors
func NewNote(noticeSpace floatgeom.Rect2, layerHeight int) int {
	layerHeight++ // Each note should be higher
	noteBaseSize := 10 + rand.Intn(5)
	noteHeight := noteBaseSize + rand.Intn(10)

	//	noteContentColor := color.RGBA{uint8(rand.Intn(220)), uint8(rand.Intn(220)), uint8(rand.Intn(220)), 255}

	var noteContentColor color.RGBA

	for k, v := range colornames.Map {
		if !strings.Contains(k, "brown") && !strings.Contains(k, "gray") {
			noteContentColor = v
		}
		break
	}
	// noteContentColor := v

	c := render.NewColorBox(noteBaseSize, noteHeight, noteContentColor)

	noteXOff := noticeSpace.W() - float64(noteBaseSize)
	noteYOff := noticeSpace.H() - float64(noteHeight)

	noteLocX := noticeSpace.Min.X() + float64(rand.Intn(int(math.Max(1, noteXOff))))
	noteLocY := noticeSpace.Min.Y() + float64(rand.Intn(int(math.Max(1, noteYOff))))

	c.SetPos(noteLocX, noteLocY)
	// c.SetPos(noticeSpace.Min.X(), noticeSpace.Min.Y())

	render.Draw(c, 2, layerHeight)
	scrawlDistance := float64(5)
	scrawlOffset := floatgeom.Point2{
		0,
		1,
	}
	for i := 0; scrawlOffset.Y() < float64(noteHeight-4); i++ {
		layerHeight++
		scrawlOffset[0] = float64(2 + rand.Intn(noteBaseSize-9))
		scrawlOffset[1] += float64(rand.Intn(noteHeight - int(scrawlOffset.Y()) - 2))

		l := render.NewThickLine(noteLocX+scrawlOffset.X(), noteLocY+scrawlOffset.Y(),
			noteLocX+scrawlOffset.X()+scrawlDistance, noteLocY+scrawlOffset.Y(),
			color.RGBA{10, 10, 10, 255}, 1)
		render.Draw(l, 2, layerHeight)
		scrawlOffset[1] += 2 + float64(rand.Intn(noteHeight-int(scrawlOffset.Y())))
	}

	return layerHeight
}
