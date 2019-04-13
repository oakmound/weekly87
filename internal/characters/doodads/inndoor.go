package doodads

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

type InnDoor struct {
	*entities.Solid
	Unmoving
}

func (d *InnDoor) Init() event.CID {
	return event.NextID(d)
}

func (d *InnDoor) Activate() {}

func NewInnDoor() *InnDoor {

	width := 83.0
	height := 258.0

	d := &InnDoor{}
	d.Solid = entities.NewSolid(float64(oak.ScreenWidth)-width, 239, width, height, nil, nil, d.Init())
	d.UpdateLabel(collision.Label(labels.Door))
	return d
}