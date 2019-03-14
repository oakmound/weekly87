package doodads

import (
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

type OutDoor struct {
	*entities.Reactive
	Unmoving
}

func (d *OutDoor) Init() event.CID {
	return event.NextID(d)
}

func NewOutDoor(runback bool) *OutDoor {

	width := 10.0
	height := float64(oak.ScreenHeight * 2 / 3)

	d := &OutDoor{}

	asset, _ := render.LoadSprite("", filepath.Join("raw", "goal.png"))

	d.Reactive = entities.NewReactive(0, float64(oak.ScreenHeight/3), width, height, asset, nil, d.Init())

	// d.Reactive = entities.NewReactive(0, 0, width, height, render.NewColorBox(int(width), int(height), color.RGBA{0, 0, 255, 255}), nil, d.Init())
	if runback {
		d.RSpace.UpdateLabel(labels.Door)
	}
	return d
}
