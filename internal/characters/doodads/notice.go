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

	// Scrawls should have the following:
	// 1) be inside the note
	// 2) be close to the last scrawl if exists
	// 3) follow a path that has gaps and heigh changes
	scrawlDistance := float64(5)
	scrawlOffset := floatgeom.Point2{
		0,
		2, //1 to in and 1 for the fact that we can offset up to 2 above
	}
	scrawlOffset[1] += float64(rand.Intn(noteHeight - int(scrawlOffset.Y()) - 2))
	for i := 0; scrawlOffset.Y() < float64(noteHeight-4); i++ {
		layerHeight++
		scrawlOffset[0] = float64(2 + rand.Intn(noteBaseSize-9))

		l := render.NewThickLine(noteLocX+scrawlOffset.X(), noteLocY+scrawlOffset.Y(),
			noteLocX+scrawlOffset.X()+scrawlDistance, noteLocY+scrawlOffset.Y(),
			color.RGBA{10, 10, 10, 255}, 1)
		render.Draw(l, 2, layerHeight)
		scrawlOffset[1] += 2
		if rand.Float64() < .5 {
			break
		}
	}

	return layerHeight
}

/*
func NewThickLine(x1, y1, x2, y2 float64, c color.Color, thickness int) *render.Sprite {
	colorer := render.IdentityColorer(c)
	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	rgba = dlb(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), colorer, thickness)
	return render.NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}

func dlb(x1, y1, x2, y2 int, colorer render.Colorer, thickness int) *image.RGBA {

	// Bresenham's line-drawing algorithm from wikipedia
	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	if xDelta == 0 && yDelta == 0 {
		width := 1 + 2*thickness
		rect := image.Rect(0, 0, width, width)
		rgba := image.NewRGBA(rect)
		for xm := 0; xm < width; xm++ {
			for ym := 0; ym < width; ym++ {
				rgba.Set(xm, ym, colorer(1.0))
			}
		}
		return rgba
	} else if xDelta == 0 {
		width := 1 + 2*thickness
		height := int(math.Floor(yDelta)) + 2*thickness
		rect := image.Rect(0, 0, width, height)
		rgba := image.NewRGBA(rect)
		for xm := 0; xm < width; xm++ {
			for ym := 0; ym < height; ym++ {
				rgba.Set(xm, ym, colorer(float64(ym)/float64(height)))
			}
		}
		return rgba
	}

	// Todo: document why we add one here
	// It has something to do with zero-height rgbas, but is always useful
	h := int(yDelta) + 1

	rect := image.Rect(0, 0, int(xDelta)+2*thickness, h+2*thickness)
	rgba := image.NewRGBA(rect)

	x2 += thickness
	y2 += thickness
	x1 += thickness
	y1 += thickness

	DrawLineColored(rgba, x1, y1, x2, y2, thickness, colorer)

	return rgba
}
*/
