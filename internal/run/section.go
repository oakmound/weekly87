package run

import (
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

var (
	baseSeed       int64
	currentSection int64
)

// Todo generate section

type Section struct {
	id       int64
	ground   *render.Sprite
	wall     *render.Sprite
	entities []characters.Character
}

func (s *Section) Draw() {
	render.Draw(s.ground, 0)
	render.Draw(s.wall, 1)
	for _, e := range s.entities {
		render.Draw(e.GetRenderable(), 2, 1)
	}
}

func (s *Section) Shift(shift float64) {
	s.ground.ShiftX(shift)
	for _, e := range s.entities {
		ShiftMoverX(e, shift)
	}
}

func (s *Section) SetBackgroundX(x float64) {
	s.ground.SetX(x)
}

func (s *Section) Destroy() {
	s.ground.Undraw()
	s.wall.Undraw()
	for _, e := range s.entities {
		e.Destroy()
	}
}

// X returns the leftmost x value of this section
func (s *Section) X() float64 {
	return s.ground.X()
}

// W returns how wide this section is
func (s *Section) W() float64 {
	// assumes ground and wall are same width
	w, _ := s.ground.GetDims()
	return float64(w)
}
