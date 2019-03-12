package section

import "github.com/oakmound/oak/dlog"

type ChangeType int

const (
	EntityDestroyed ChangeType = iota
)

type Change struct {
	typ ChangeType
	val int
}

func (s *Section) ApplyChange(ch *Change) {
	switch ch.typ {
	case EntityDestroyed:
		// val is index of entity destroyed
		if ch.val > len(s.entities) {
			dlog.Error("Entity to destroy", ch.val, "does not exist in section")
			return
		}
		s.entities[ch.val] = nil
	default:
		dlog.Error("Unknown section change type:", ch.typ)
	}
}
