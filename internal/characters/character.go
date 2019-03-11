package characters

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities/x/move"
)

type Character interface {
	move.Mover
	Destroy()
	GetReactiveSpace() *collision.ReactiveSpace
}

type Player interface {
	Character
	Special1()
}

const (
	playerHeight = 32
	playerWidth  = 16
)
