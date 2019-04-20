package records

import (
	"encoding/json"
	"math/rand"
	"os"

	"github.com/oakmound/oak/dlog"
)

const recordsFile = "save.json"

// Records serves as our safe file and all variables we track across
// multiple runs
type Records struct {
	SectionsCleared int64 `json:"sectionsCleared"`
	BaseSeed        int64 `json:"baseSeed"`
	// Todo: more
	FarthestGoneInSections int64 `json:"farthestGoneInSections"`
}

func (s *Records) Store() {
	f, err := os.Create(recordsFile)
	data, err := json.Marshal(s)
	dlog.ErrorCheck(err)
	_, err = f.Write(data)
	dlog.ErrorCheck(err)
	dlog.ErrorCheck(f.Close())
}

func Load() *Records {
	r := &Records{}

	f, err := os.Open(recordsFile)
	if err != nil {
		f, err := os.Create(recordsFile)
		dlog.ErrorCheck(err)
		r.BaseSeed = rand.Int63()
		data, err := json.Marshal(r)
		dlog.ErrorCheck(err)
		_, err = f.Write(data)
		dlog.ErrorCheck(err)
	} else {
		dc := json.NewDecoder(f)
		dlog.ErrorCheck(dc.Decode(r))
	}
	if f != nil {
		dlog.ErrorCheck(f.Close())
	}

	return r
}
