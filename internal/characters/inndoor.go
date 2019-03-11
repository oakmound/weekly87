package characters

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

// Door is a small struct to make Doors (initially just the 2 innDoors)
type Door struct {
	*entities.Solid
	Unmoving
}

func (d *Door) Init() event.CID {
	return event.NextID(d)
}

func NewDoor() *Door {

	width := float64(oak.ScreenWidth / 8)
	height := float64(oak.ScreenHeight)

	d := &Door{}
	d.Solid = entities.NewSolid(0, 0, width, height, render.NewColorBox(int(width), int(height), color.RGBA{0, 0, 255, 255}), nil, d.Init())
	d.UpdateLabel(collision.Label(LabelDoor))
	return d
}

func (d *Door) Flip() {
}
