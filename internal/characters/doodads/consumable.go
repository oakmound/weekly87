package doodads

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/layer"
)

type Consumer interface {
	GetPos()
}

type Consumable struct {
	*entities.Solid
	R render.Modifiable
}

// Init creates the furnitures entity id
func (c *Consumable) Init() event.CID {
	return event.NextID(c)
}

// NewConsumable creates a consumable from an image
func NewConsumable(x, y float64, img *render.Sprite) *Consumable {
	c := &Consumable{}
	c.R = img.Copy()
	width, height := img.GetDims()
	c.Solid = entities.NewSolid(x, y, float64(width), float64(height), nil, nil, c.Init())
	c.UpdateLabel(collision.Label(labels.Drinkable))
	c.SetPos(x, y)
	render.Draw(c.R, layer.Play, 2)

	return c
}

// SetPos of the solid and the renderable
func (c *Consumable) SetPos(x, y float64) {
	c.Solid.SetPos(x, y)
	c.R.SetPos(x, y)
}

func (c *Consumable) Consume() {

}
