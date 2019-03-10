package characters

import (
	"github.com/oakmound/oak/entities/x/move"
)

type Character interface {
	move.Mover
	Destroy()
}

type Player interface {
	Character
	Attack1()
}

const (
	playerHeight = 32
	playerWidth  = 16
)
