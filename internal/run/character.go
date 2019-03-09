package run

import (
	"github.com/oakmound/oak/entities/x/move"
)

type Character interface {
	move.Mover
	// ???
}

type Player interface {
	Character
	Attack1()
}

const (
	playerHeight = 32
	playerWidth  = 16
)
