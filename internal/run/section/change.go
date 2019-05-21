package section

import (
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/weekly87/internal/characters"
)

type ChangeType int

const (
	EntityDestroyed ChangeType = iota
	EntityAdded
)

type Change struct {
	Typ    ChangeType
	Val    int
	Entity characters.Character
}

func (s *Section) ApplyChange(ch Change) {
	switch ch.Typ {
	case EntityDestroyed:
		// val is index of entity destroyed
		if ch.Val > len(s.entities) {
			dlog.Error("Entity to destroy", ch.Val, "does not exist in section")
			return
		}
		s.entityMutex.Lock()
		s.entities[ch.Val] = nil
		s.entityMutex.Unlock()
	case EntityAdded:
		s.entityMutex.Lock()
		s.entities = append(s.entities, ch.Entity)
		s.entityMutex.Unlock()
	default:
		dlog.Error("Unknown section change type:", ch.Typ)
	}
}
