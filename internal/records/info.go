package records

import "github.com/oakmound/weekly87/internal/characters/players"

type RunInfo struct {
	Party           []*players.Player
	SectionsCleared int
}
