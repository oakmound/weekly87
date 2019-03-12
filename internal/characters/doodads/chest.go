package doodads

import (
	"image/color"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

type Chest struct {
	*entities.Reactive
	Unmoving
	Value int64
}

func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

func NewChest(value int64) *Chest {
	ch := &Chest{}
	// Todo: calculate image based on value
	ch.Reactive = entities.NewReactive(0, 0, 16, 16,
		render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255}), nil, ch.Init())
	ch.RSpace.UpdateLabel(labels.Chest)
	ch.Value = value
	return ch
}
