package run

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type Spearman struct {
	*entities.Moving
}

func (s *Spearman) Init() event.CID {
	return event.NextID(s)
}

func NewSpearman(x, y float64) *Spearman {
	s := &Spearman{}
	r := render.NewColorBox(playerWidth, playerHeight, color.RGBA{255, 0, 0, 255})
	s.Moving = entities.NewMoving(x, y, playerWidth, playerHeight, r, nil, s.Init(), 0)
	s.Speed = physics.NewVector(0, 5)
	return s
}

func (s *Spearman) Attack1() {
	fmt.Println("Attacking!")
}
