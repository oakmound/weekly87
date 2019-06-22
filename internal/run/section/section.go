package section

import (
	"sync"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/entities/x/move"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/layer"
)

var (
	baseSeed       int64
	currentSection int64
)

type Section struct {
	id          int64
	ground      *render.Sprite
	wall        *render.Sprite
	entities    []characters.Character
	entityMutex sync.Mutex
}
type MoverWithParticles interface {
	MoveParticles(floatgeom.Point2)
}

func (s *Section) Copy() *Section {
	return &Section{
		id:          s.id,
		ground:      s.ground.Copy().(*render.Sprite),
		wall:        s.wall.Copy().(*render.Sprite),
		entityMutex: sync.Mutex{},
	}
}

func (s *Section) Draw() {
	render.Draw(s.ground, layer.Ground)
	render.Draw(s.wall, layer.Background)
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
		if e != nil {
			e.Destroy()
		}
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
		if e != nil {
			e.Activate()
			render.Draw(e.GetRenderable(), layer.Play, 1)
		}
	}
}

func (s *Section) ShiftEntities(shift float64) {
	s.entityMutex.Lock()
	for _, e := range s.entities {
		if e != nil {
			move.ShiftX(e, shift)
			if pm, ok := e.(MoverWithParticles); ok {
				pm.MoveParticles(floatgeom.Point2{shift, 0})
			}
		}
	}
	s.entityMutex.Unlock()
}

func (s *Section) AppendEntities(e ...characters.Character) {
	s.entityMutex.Lock()
	s.entities = append(s.entities, e...)
	s.entityMutex.Unlock()
}

func (s *Section) GetId() int64 {
	return s.id
}
