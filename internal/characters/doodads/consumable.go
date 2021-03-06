package doodads

import (
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/timing"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/layer"
)

// Consumer can consume consumables... TODO: add to this comment as we update this interface
type Consumer interface {
	GetPos()
}

// Consumable is an object that can be drawn, has collision and can be consumed
type Consumable struct {
	*entities.Solid
	R      render.Modifiable
	cLabel collision.Label
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

	c.SetPos(x, y)
	render.Draw(c.R, layer.Play, 2)

	return c
}

// Activate the consumable allowing it to be consumed
// Consider: Should this report on something? (already activated? or failed to activate?)
func (c *Consumable) Activate() {
	c.Bind(func(id int, space interface{}) int {
		s := space.(*collision.Space)
		c.Consume(0, 20, s)
		return 1
	}, "Consume")
	c.UpdateLabel(collision.Label(c.cLabel))
	return
}

// SetPos of the solid and the renderable for the consumable
func (c *Consumable) SetPos(x, y float64) {
	c.Solid.SetPos(x, y)
	c.R.SetPos(x, y)
}

// ShiftPos of the solid and the renderable
func (c *Consumable) ShiftPos(x, y float64) {
	c.Solid.ShiftPos(x, y)
	c.R.ShiftX(x)
	c.R.ShiftY(y)
}

// Consume in what is for now a boiler plate fashion
func (c *Consumable) Consume(xOff, yOff float64, sp *collision.Space) {
	dlog.Info("pretending to be consumed")
	c.Tree.Remove(c.Space)
	x := sp.X() - xOff
	y := sp.Y() - yOff

	speedMod := .25

	delta := physics.NewVector(x-c.X(), y-c.Y()).Normalize().Scale(speedMod)

	c.Bind(func(id int, _ interface{}) int {

		con, ok := event.GetEntity(id).(*Consumable)
		if !ok {
			dlog.Error("Non Consumer in consume enter")
		}

		con.ShiftPos(delta.X(), delta.Y())

		if curX, curY := con.GetPos(); curX > x && curY < y {
			c.R.Undraw()
			c.Destroy()

			// Let the consumer know that it is no longer consuming
			sp.CID.Trigger("consumeCompleted", con.CID)
			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")

}

// NewDrinkable creates a drinkable consumable
func NewDrinkable(x, y float64, img *render.Sprite) *Consumable {
	c := NewConsumable(x, y, img)
	c.cLabel = labels.Drinkable
	go timing.DoAfter(40 * time.Second, func() {
		c.cLabel = 0
		c.Destroy()
		c.R.Undraw()
	})
	return c
}
