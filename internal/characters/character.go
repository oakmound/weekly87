package characters

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
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
	Alive() bool
	Kill()
}

type basePlayer struct {
	*entities.Interactive
	alive bool
}

func (bp *basePlayer) Alive() bool {
	return bp.alive
}

func (bp *basePlayer) Kill() {
	bp.alive = false
	// Todo: animation
}

const (
	playerHeight = 32
	playerWidth  = 16
)
