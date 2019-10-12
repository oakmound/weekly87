package keyviz

import (
	"image"
	"image/color"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
)

func New(opts ...Option) image.Image {
	g := &Generator{}
	for _, o := range opts {
		o(g)
	}
	return g.Generate()
}

type Generator struct {
	Text     string
	TextSize float64
	// These will be inferred based on TextSize,
	// and vice versa, if not provided
	Width        float64
	Height       float64
	FontFile     string
	CornerOffset float64
	TextXOffset  float64
	TextYOffset  float64
	// Style TODO
	Color     color.Color
	LineColor color.Color
}

func (g Generator) Generate() image.Image {
	font := render.DefFont()
	if g.FontFile != "" {
		font.File = g.FontFile
	}
	if g.TextSize != 0 {
		font.Size = g.TextSize
	}
	font = font.Copy()
	txt := font.NewStrText(g.Text, 0, 0)
	w, h := txt.GetDims()
	if g.CornerOffset == 0 {
		g.CornerOffset = 7
	}
	if g.TextXOffset == 0 {
		g.TextXOffset = 2
	}
	if g.TextYOffset == 0 {
		g.TextYOffset = 1
	}
	if g.Width == 0 {
		// infer width based on text and text size
		g.Width = float64(w) + g.CornerOffset*2 + g.TextXOffset*2
	}
	if g.Height == 0 {
		// infer height based on text size
		g.Height = float64(h) + g.CornerOffset*3 + g.TextYOffset*2
	}
	if g.LineColor == nil {
		g.LineColor = color.RGBA{50, 50, 50, 200}
	}

	// .         .
	//  . . . . .
	//  . text  .
	//  . . . . .
	// .         .
	//.           .

	rect := floatgeom.NewRect2WH(g.CornerOffset, g.CornerOffset, float64(w)+g.TextXOffset*2, float64(h)+g.TextYOffset*2)

	txtSpr := txt.ToSprite()
	txtSpr.SetPos(rect.Min.X()+g.TextXOffset, rect.Min.Y()+g.TextYOffset)

	comp := render.NewCompositeM(
		render.NewColorBox(int(g.Width), int(g.Height), g.Color),
		// Corners
		render.NewThickLine(0, 0, rect.Min.X(), rect.Min.Y(), g.LineColor, 1),
		render.NewThickLine(g.Width, 0, rect.Max.X(), rect.Min.Y(), g.LineColor, 1),
		render.NewThickLine(g.Width, g.Height, rect.Max.X(), rect.Max.Y(), g.LineColor, 1),
		render.NewThickLine(0, g.Height, rect.Min.X(), rect.Max.Y(), g.LineColor, 1),
		// Sides
		render.NewThickLine(rect.Min.X(), rect.Min.Y(), rect.Min.X(), rect.Max.Y(), g.LineColor, 1),
		render.NewThickLine(rect.Min.X(), rect.Max.Y(), rect.Max.X(), rect.Max.Y(), g.LineColor, 1),
		render.NewThickLine(rect.Max.X(), rect.Max.Y(), rect.Max.X(), rect.Min.Y(), g.LineColor, 1),
		render.NewThickLine(rect.Max.X(), rect.Min.Y(), rect.Min.X(), rect.Min.Y(), g.LineColor, 1),

		txtSpr,
	)

	out := comp.ToSprite().GetRGBA()
	return out
}

type Option func(g *Generator)
