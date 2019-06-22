package players

import (
	"github.com/oakmound/weekly87/internal/abilities/buff"
)

// A Buffer is an abilty that bestow a set of buffs
type Buffer interface {
	Buffs() []buff.Buff
}

// A Destroyable can be destroyed
type Destroyable interface {
	Destroy()
}
