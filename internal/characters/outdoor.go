package characters

import (
	"image/color"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type OutDoor struct {
	*entities.Reactive
	Unmoving
}

func (d *OutDoor) Init() event.CID {
	return event.NextID(d)
}

func NewOutDoor(runback bool) *OutDoor {

	width := float64(oak.ScreenWidth / 8)
	height := float64(oak.ScreenHeight)

	d := &OutDoor{}
	d.Reactive = entities.NewReactive(0, 0, width, height, render.NewColorBox(int(width), int(height), color.RGBA{0, 0, 255, 255}), nil, d.Init())
	if runback {
		d.RSpace.UpdateLabel(LabelDoor)
	}
	return d
}
