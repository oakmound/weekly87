package section

import (
	"github.com/oakmound/oak/entities/x/move"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

var (
	baseSeed       int64
	currentSection int64
)

type Section struct {
	id       int64
	ground   *render.Sprite
	wall     *render.Sprite
	entities []characters.Character
}

func (s *Section) Copy() *Section {
	return &Section{
		id:     s.id,
		ground: s.ground.Copy().(*render.Sprite),
		wall:   s.wall.Copy().(*render.Sprite),
	}
}

func (s *Section) Draw() {
	render.Draw(s.ground, 0)
	render.Draw(s.wall, 1)
}

func (s *Section) Shift(shift float64) {
	s.wall.ShiftX(shift)
	s.ground.ShiftX(shift)
	s.ShiftEntities(shift)
}

func (s *Section) SetBackgroundX(x float64) {
	delta := x - s.wall.X()
	s.wall.SetX(x)
	s.ground.SetX(x)
	s.ShiftEntities(delta)
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
	return s.wall.X()
}

// W returns how wide this section is
func (s *Section) W() float64 {
	// assumes ground and wall are same width
	w, _ := s.wall.GetDims()
	return float64(w)
}

func (s *Section) ActivateEntities() {
	for _, e := range s.entities {
		e.Activate()
		render.Draw(e.GetRenderable(), 2, 1)
	}
}

func (s *Section) ShiftEntities(shift float64) {
	for _, e := range s.entities {
		move.ShiftX(e, shift)
	}
}
