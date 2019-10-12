package labels

import (
	"image/color"

	"github.com/oakmound/oak/collision"
)

const (
	None = iota
	Door
	Chest
	PC
	Enemy
	NPC
	Blocking
	Drinkable
	Ornament
	EffectsPlayer
	EffectsEnemy
)

var ColorMap = map[collision.Label]color.RGBA{
	Chest:    color.RGBA{255, 255, 0, 255},
	Door:     color.RGBA{125, 125, 125, 255},
	Enemy:    color.RGBA{0, 0, 255, 255},
	PC:       color.RGBA{125, 0, 255, 255},
	Blocking: color.RGBA{200, 200, 10, 255},
	Ornament: color.RGBA{250, 200, 40, 255},
	NPC:      color.RGBA{125, 200, 10, 255},
}
