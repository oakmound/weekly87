package doodads

import (
	"path/filepath"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

type Chest struct {
	*entities.Reactive
	Unmoving
	Value int64
}

func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

func NewChest(value int64) *Chest {
	ch := &Chest{}
	// r := render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255})
	r, _ := render.LoadSprite("", filepath.Join("", "/16x16/chest.png"))
	// Todo: calculate image based on value
	ch.Reactive = entities.NewReactive(0, 0, 16, 16,
		r, nil, ch.Init())

	ch.RSpace.UpdateLabel(labels.Chest)
	ch.Value = value
	return ch
}

func (c *Chest) Destroy() {
	c.Doodad.Destroy()
	c.Tree.Remove(c.RSpace.Space)
	c.RSpace.Space.UpdateLabel(0)
}
