package characters

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities/x/move"
)

type Character interface {
	move.Mover
	Destroy()
	GetReactiveSpace() *collision.ReactiveSpace
	Activate()
}
