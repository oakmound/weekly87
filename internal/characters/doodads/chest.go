package doodads

import (
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/recolor"
	"github.com/oakmound/weekly87/internal/restrictor"
)

// Chest is a unit of value which can be interacted with and drawn
type Chest struct {
	*entities.Reactive
	Unmoving
	Value  int64
	Active bool
}

// Init the chest and get its CID
func (c *Chest) Init() event.CID {
	return event.NextID(c)
}

// Destroy a chest and clean up its artifcats
func (c *Chest) Destroy() {
	c.Active = false
	c.Reactive.Destroy()
}

// Activate a Chest and allow users to interact with it
func (c *Chest) Activate() {
	restrictor.Add(c)
	c.Active = true
}

// GetDims of the chest renderable
func (c *Chest) GetDims() (int, int) {
	return c.Reactive.R.GetDims()
}

// NewChest creates a chest with the given value
func NewChest(value int64) *Chest {
	ch := &Chest{}
	// r := render.NewColorBox(16, 16, color.RGBA{0, 255, 255, 255})
	r, _ := render.LoadSprite("", filepath.Join("", "/16x16/chest.png"))
	r = r.Copy().(*render.Sprite)
	w, h := 16.0, 16.0
	switch value {
	case 1:
		break
	case 2:
		// recolor blue
		r.Filter(recolor.WithStrategy(recolor.ColorShift(color.RGBA{100, 100, 255, 10})))
	case 3:
		// recolor purple
		r.Filter(recolor.WithStrategy(recolor.ColorShift(color.RGBA{255, 100, 255, 10})))
	case 4:
		// recolor red
		r.Filter(recolor.WithStrategy(recolor.ColorShift(color.RGBA{255, 100, 100, 10})))
	case 5:
		// size up
		r.Modify(mod.Scale(2, 2))
		w *= 2
		h *= 2
	}
	// Todo: calculate image based on value
	ch.Reactive = entities.NewReactive(0, 0, w, h,
		r, nil, ch.Init())

	ch.RSpace.UpdateLabel(labels.Chest)
	ch.Value = value
	return ch
}
