package characters

import (
	"github.com/oakmound/oak/entities/x/move"
)

type Character interface {
	move.Mover
	// ???
}

type Player interface {
	Character
	Special1()
	Special2() //Should Return an ability
}

const (
	playerHeight = 32
	playerWidth  = 16
)
