package doodads

import (
	"path/filepath"

	"github.com/oakmound/oak/dlog"

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

	tiles, err := render.LoadSprites(filepath.Join("assets", "images"),
		filepath.Join("raw", "goal.png"), 16, 400, 0)
	dlog.ErrorCheck(err)

	swtch := render.NewSwitch(
		"uncut",
		map[string]render.Modifiable{
			"uncut": tiles[0][0],
			"cut":   tiles[1][0],
		},
	)

	d.Reactive = entities.NewReactive(0, 0, width, height, swtch, nil, d.Init())

	d.Bind(func(id int, _ interface{}) int {
		dr, ok := event.GetEntity(id).(*OutDoor)
		if ok {
			dr.R.(*render.Switch).Set("cut")
		} else {
			return event.UnbindSingle
		}
		return 0
	}, "RibbonCut")

	if runback {
		d.RSpace.UpdateLabel(labels.Door)
	}
	return d
}
