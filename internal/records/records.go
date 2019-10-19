package records

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/oakmound/weekly87/internal/characters/players"

	"github.com/oakmound/oak/dlog"
)

const recordsFile = "save.json"
const archPath = "save_arch_"

// Records serves as our save file and all variables we track across
// multiple runs
type Records struct {
	SectionsCleared int64 `json:"sectionsCleared"`
	BaseSeed        int64 `json:"baseSeed"`
	// Todo: more
	FarthestGoneInSections int64                 `json:"farthestGoneInSections"`
	EnemiesDefeated        int64                 `json:"enemiesDefeated"`
	PartyComp              []players.PartyMember `json:"partyComp"`
	Deaths                 int                   `json:"deaths"`
	Wealth                 int                   `json:"wealth"`

	LastRun RunInfo `json:"lastRun"`
}

// Store a record to a file
func (s *Records) Store() {
	f, err := os.Create(recordsFile)
	data, err := json.Marshal(s)
	dlog.ErrorCheck(err)
	_, err = f.Write(data)
	dlog.ErrorCheck(err)
	dlog.ErrorCheck(f.Close())
}

// Load a record from a file
func Load() *Records {
	r := &Records{}

	f, err := os.Open(recordsFile)
	if err != nil {
		f, err := os.Create(recordsFile)
		dlog.ErrorCheck(err)
		r.BaseSeed = rand.Int63()
		r.PartyComp = []players.PartyMember{{PlayerClass: players.Swordsman, AccruedValue: 0, Name: "Dan the Default"}}
		r.LastRun = RunInfo{EnemiesDefeated: 0, SectionsCleared: 0}
		data, err := json.Marshal(r)
		dlog.ErrorCheck(err)
		_, err = f.Write(data)
		dlog.ErrorCheck(err)
	} else {
		dc := json.NewDecoder(f)
		dlog.ErrorCheck(dc.Decode(r))
		if r.PartyComp == nil {
			r.PartyComp = []players.PartyMember{{PlayerClass: players.Swordsman, AccruedValue: 0, Name: "Dan the Default"}}
		}
	}
	if f != nil {
		dlog.ErrorCheck(f.Close())
	}

	return r
}

// Archive the current save file
func Archive() (string, error) {
	newName := fmt.Sprintf("%s%s.json", archPath, time.Now().Format("MonJan2150405"))
	err := os.Rename(recordsFile, newName)
	return newName, err
}
