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

func (r *Records) Store() {
	f, err := os.Open("save.json")
	if err != nil {
		f, _ = os.Create("save.json")
	}
	dc := json.NewDecoder(f)
	dlog.ErrorCheck(dc.Decode(r))
}
func Load() *Records {
	r := &Records{}

	f, err := os.Open("save.json")
	if err != nil {
		f, err := os.Create("save.json")
		if err != nil {
			dlog.Error(err)
		}
		r.BaseSeed = rand.Int63()
		data, err := json.Marshal(r)
		if err != nil {
			dlog.ErrorCheck(err)
		}
		f.Write(data)
	} else {
		dc := json.NewDecoder(f)
		dlog.ErrorCheck(dc.Decode(r))
	}
	return r
}
