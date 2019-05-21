package doodads

import (
	"math/rand"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

// Ornament creates an ornament randomly within the given space trying not to overlap
type Ornament struct {
	*entities.Solid
	Unmoving
	R render.Modifiable
}

// Init creates the furnitures entity id
func (o *Ornament) Init() event.CID {
	return event.NextID(o)
}

// Activate the ornament, doing nothing
func (o *Ornament) Activate() {}

// NewOrnament creates a new ornament and places it if possible
// This should only be run after creating static ornaments
func NewOrnament(x, y, w, h float64, img *render.Sprite) *Ornament {
	o := &Ornament{}
	o.R = img.Copy()
	// Decide the location
	width, height := img.GetDims()
	wOffset := 0.0
	hOffset := 0.0

	collided := true

	maxTries := 20
	o.Solid = entities.NewSolid(x, y, float64(width), float64(height), nil, nil, o.Init())
	for tries := 0; collided; tries++ {
		if tries == maxTries {
			o.Tree.Add(o.Space)
			return nil
		}
		collided = false
		wOffset = float64(rand.Intn(int(w) - width))
		hOffset = float64(rand.Intn(int(h) - height))
		o.Solid.SetPos(x+wOffset, y+hOffset)

		o.Space.Label = collision.Label(labels.Ornament)

		if overlapSpace := collision.HitLabel(o.Space, collision.Label(labels.Ornament)); overlapSpace != nil {
			collided = true
		}

	}

	o.R.SetPos(x+wOffset, y+hOffset)
	render.Draw(o.R, 2, 3)
	return o
}
