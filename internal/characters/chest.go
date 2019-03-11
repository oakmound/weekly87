package characters

import (
	"image/color"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var _ Character = &Chest{}

type Chest struct {
	*entities.Reactive
	Unmoving
}

func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

func NewChest(value int64) *Chest {
	ch := &Chest{}
	ch.Reactive = entities.NewReactive(0, 0, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255}), nil, ch.Init())
	ch.RSpace.UpdateLabel(LabelChest)
	return ch
}
