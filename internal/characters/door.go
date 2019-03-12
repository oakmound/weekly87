package characters

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
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

	width := 83.0
	height := 258.0

	d := &Door{}
	d.Solid = entities.NewSolid(float64(oak.ScreenWidth)-width, 239, width, height, nil, nil, d.Init())
	d.UpdateLabel(collision.Label(LabelDoor))
	return d
}

func (d *Door) Flip() {
}
