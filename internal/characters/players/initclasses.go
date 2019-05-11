package players

import (
	"strings"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

func Init() {
	WarriorsInit()
	MageInit()
	classmapping = map[int]*Constructor{
		Spearman:  SpearmanConstructor,
		Swordsman: SwordsmanConstructor,
		Mage:      MageConstructor,
	}
}
func filterCharMap(baseCharMap map[string]render.Modifiable, filter mod.Filter) map[string]render.Modifiable {
	outputMap := make(map[string]render.Modifiable)

	for k, v := range baseCharMap {
		outputMap[k] = v.Copy()
		if !strings.Contains(k, "dead") {
			outputMap[k].Filter(filter)
		}
	}

	return outputMap
}

const (
	Spearman = iota
	Swordsman
	Mage
)

var classmapping map[int]*Constructor

// ClassConstructor creates the classes from a int64 list
func ClassConstructor(partyComp []int) []Constructor {
	classes := make([]Constructor, len(partyComp))
	for i, c := range partyComp {
		classes[i] = *classmapping[c].Copy()
	}
	return classes
}
