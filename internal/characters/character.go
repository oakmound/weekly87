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

const (
	playerHeight = 32
	playerWidth  = 16
)
