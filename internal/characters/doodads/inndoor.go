package doodads

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

// InnDoor is for usage in the inn scene
type InnDoor struct {
	*entities.Solid
	Unmoving
	NextScene string
}

func (d *InnDoor) Init() event.CID {
	return event.NextID(d)
}

func (d *InnDoor) Activate() {}

// NewInnDoor creates a door from the inn scene
func NewInnDoor(nextscene string) *InnDoor {
	width := 83.0
	height := 258.0
	return NewCustomInnDoor(nextscene, float64(oak.ScreenWidth)-width, 239, width, height)
}

func NewCustomInnDoor(nextscene string, x, y, w, h float64) *InnDoor {
	d := &InnDoor{}
	d.Solid = entities.NewSolid(x, y, w, h, nil, nil, d.Init())
	d.UpdateLabel(collision.Label(labels.Door))
	d.NextScene = nextscene
	return d
}
