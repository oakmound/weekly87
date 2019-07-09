package inn

import (
	"errors"
	"image/color"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/render"
)

type SelectorConstructor struct {
	Start    floatgeom.Point2
	Size     intgeom.Point2
	Display  func(size intgeom.Point2) render.Renderable
	Step     floatgeom.Point2
	Limit    int
	Layers   []int
	Callback func(i int)
}

type Selector struct {
	*SelectorConstructor
	pos int
	R   render.Renderable
}

func (sc *SelectorConstructor) Generate() (*Selector, error) {
	if sc.Display == nil {
		sc.Display = func(size intgeom.Point2) render.Renderable {
			return render.NewColorBox(size.X(), size.Y(), color.RGBA{125, 125, 0, 125})
		}
	}

	s := &Selector{
		SelectorConstructor: sc,
		R:                   sc.Display(sc.Size),
	}

	render.Draw(s.R, sc.Layers...)

	return s, nil
}

func (s *Selector) MoveTo(i int) error {
	if i >= s.Limit {
		return errors.New("Index to move to exceeds limit")
	}
	delta := float64(i - s.pos)
	s.R.ShiftX(s.Step.X() * delta)
	s.R.ShiftY(s.Step.Y() * delta)
	s.pos = i
	return nil
}

func (s *Selector) Select() {
	s.Callback(s.pos)
}
