package doodads

import (
	"math/rand"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/labels"
)

// Furniture is made for elements in the inn that the player should not be able to walk through.
type Furniture struct {
	*entities.Solid
	Unmoving
}

// Init creates the furnitures entity id
func (f *Furniture) Init() event.CID {
	return event.NextID(f)
}

// Activate the furniture, doing nothing
func (f *Furniture) Activate() {}

// NewFurniture creates a new piece of Furniture
func NewFurniture(x, y, w, h float64) *Furniture {

	f := &Furniture{}
	f.Solid = entities.NewSolid(x, y, w, h, nil, nil, f.Init())
	f.UpdateLabel(collision.Label(labels.Blocking))
	return f
}

// SetOrnaments tries to place a number of a given ornament on the piece of furniture`
func (f *Furniture) SetOrnaments(possibleSprites []*render.Sprite, instanceCount int) {
	for i := 0; i < instanceCount; i++ {
		NewOrnament(f.X(), f.Y(), f.W, f.H, possibleSprites[rand.Intn(len(possibleSprites))])
	}

}
