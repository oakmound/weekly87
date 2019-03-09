package run

import (
	"image/color"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
)

type Chest struct {
	*entities.Solid
}

func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

func NewChest(value int64) *Chest {
	ch := &Chest{}
	ch.Solid = entities.NewSolid(0, 0, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255}), nil, ch.Init())
	ch.UpdateLabel(LabelChest)
	return ch
}

func (c *Chest) GetSpeed() physics.Vector {
	return physics.Vector{}
}

func (c *Chest) GetDelta() physics.Vector {
	return physics.Vector{}
}
