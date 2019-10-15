package doodads

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
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
	c.Bind(func(id int, space interface{}) int {
		s := space.(*collision.Space)
		c.Consume(s.X(), s.Y())
		return 1
	}, "Consume")

	c.SetPos(x, y)
	render.Draw(c.R, layer.Play, 2)

	return c
}

// NewDrinkable creates a drinkable consumable
func NewDrinkable(x, y float64, img *render.Sprite) *Consumable {
	c := NewConsumable(x, y, img)
	c.UpdateLabel(collision.Label(labels.Drinkable))
	return c
}

// SetPos of the solid and the renderable
func (c *Consumable) SetPos(x, y float64) {
	c.Solid.SetPos(x, y)
	c.R.SetPos(x, y)
}

// Consume in what is for now a boiler plate fashion
func (c *Consumable) Consume(x, y float64) {
	dlog.Info("pretending to be consumed")
	c.R.Undraw()
	c.Destroy()
}
