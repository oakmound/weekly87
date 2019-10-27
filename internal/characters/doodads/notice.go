package doodads

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"strings"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/shiny/materialdesign/colornames"
	"github.com/oakmound/weekly87/internal/layer"
)

// NewNote creates overlapping notes with random colors
func NewNote(noticeSpace floatgeom.Rect2, layerHeight int) int {
	layerHeight++ // Each note should be higher
	noteBaseSize := 10 + rand.Intn(8)
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

	// Background for a note via Colorybox
	c := render.NewColorBox(noteBaseSize, noteHeight, noteContentColor)

	noteXOff := (noticeSpace.W() - float64(noteBaseSize)) / 2
	noteYOff := noticeSpace.H() - float64(noteHeight)

	noteLocX := noticeSpace.Min.X() + float64(rand.Intn(int(math.Max(1, noteXOff))))
	noteLocY := noticeSpace.Min.Y() + float64(rand.Intn(int(math.Max(1, noteYOff))))

	c.SetPos(noteLocX, noteLocY)
	// c.SetPos(noticeSpace.Min.X(), noticeSpace.Min.Y())

	render.Draw(c, layer.Play, layerHeight)

	// Scrawls should have the following:
	// 1) be inside the note
	// 2) be close to the last scrawl if exists
	// 3) follow a path that has gaps and height changes
	scrawlDistance := float64(5)
	scrawlOffset := floatgeom.Point2{
		0,
		2, //1 to in and 1 for the fact that we can offset up to 2 above
	}
	scrawlOffset[1] += float64(rand.Intn((noteHeight-4)/2 - int(scrawlOffset.Y())))
	for i := 0; scrawlOffset.Y() < float64(noteHeight-4); i++ {
		layerHeight++
		scrawlOffset[0] = float64(2 + rand.Intn(noteBaseSize-9))

		l := NewPunctuatedDeviatedLine(noteLocX+scrawlOffset.X(), noteLocY+scrawlOffset.Y(),
			noteLocX+scrawlOffset.X()+scrawlDistance, noteLocY+scrawlOffset.Y(),
			render.IdentityColorer(color.RGBA{10, 10, 10, 255}), 1)
		render.Draw(l, layer.Play, layerHeight)
		scrawlOffset[1] += 4
		if rand.Float64() < .2 {
			break
		}
	}

	return layerHeight
}

//TODO: consider pulling into oak

// NewPunctuatedDeviatedLine returns a line with a custom function for how each pixel in that line should be colored.
func NewPunctuatedDeviatedLine(x1, y1, x2, y2 float64, colorer render.Colorer, thickness int) *render.Sprite {
	var rgba *image.RGBA
	// We subtract the minimum from each side here
	// to normalize the new line segment toward the origin
	minX := math.Min(x1, x2)
	minY := math.Min(y1, y2)
	rgba = drawDeviatedLineBetween(int(x1-minX), int(y1-minY), int(x2-minX), int(y2-minY), colorer, thickness)
	return render.NewSprite(minX-float64(thickness), minY-float64(thickness), rgba)
}
func drawDeviatedLineBetween(x1, y1, x2, y2 int, colorer render.Colorer, thickness int) *image.RGBA {

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

	DrawPunctuatedDeviatedLineColored(rgba, x1, y1, x2, y2, thickness, 0, 1, 0.7, colorer)

	return rgba
}

// DrawPunctuatedDeviatedLineColored has too long of a name... but hey its easier than changing oak
func DrawPunctuatedDeviatedLineColored(rgba *image.RGBA, x1, y1, x2, y2, thickness, xPixelDeviance, yPixelDeviance int, gapPercent float64, colorer render.Colorer) {

	if gapPercent >= 1 {
		gapPercent = 0.0
	}

	xDelta := math.Abs(float64(x2 - x1))
	yDelta := math.Abs(float64(y2 - y1))

	xSlope := -1
	x3 := x1
	if x2 < x1 {
		xSlope = 1
		x3 = x2
	}
	ySlope := -1
	y3 := y1
	if y2 < y1 {
		ySlope = 1
		y3 = y2
	}

	w := int(xDelta)
	h := int(yDelta)

	progress := func(x, y, w, h int) float64 {
		hprg := render.HorizontalProgress(x, y, w, h)
		vprg := render.VerticalProgress(x, y, w, h)
		if ySlope == -1 {
			vprg = 1 - vprg
		}
		if xSlope == -1 {
			hprg = 1 - hprg
		}
		return (hprg + vprg) / 2
	}

	err := xDelta - yDelta
	var err2 float64
	for i := 0; true; i++ {

		for xm := x2 - thickness - yPixelDeviance; xm <= (x2 + thickness + xPixelDeviance); xm++ {
			for ym := y2 - thickness - yPixelDeviance; ym <= (y2 + thickness + yPixelDeviance); ym++ {
				if rand.Float64() > gapPercent {
					p := progress(xm-x3, ym-y3, w, h)
					rgba.Set(xm, ym, colorer(p))
				}
			}
		}
		if x2 == x1 && y2 == y1 {
			break
		}
		err2 = 2 * err
		if err2 > -1*yDelta {
			err -= yDelta
			x2 += xSlope
		}
		if err2 < xDelta {
			err += xDelta
			y2 += ySlope
		}
	}
}
