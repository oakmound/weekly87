package players

import (
	"image/color"
	"strings"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/abilities"
)

func Init() {
	WarriorsInit()
	MageInit()
	classmapping = map[int]*Constructor{
		Spearman:  WarriorConstructors["Spearman"],
		Swordsman: WarriorConstructors["Swordsman"],
		Berserker: WarriorConstructors["Berserker"],
		Paladin:   WarriorConstructors["Paladin"],
		Mage:      MageConstructors["Red"],
		WhiteMage: MageConstructors["White"],
		BlueMage:  MageConstructors["Blue"],
		TimeMage:  MageConstructors["Time"],
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
	Swordsman = iota
	Berserker
	Paladin
	Spearman
	Mage
	WhiteMage
	BlueMage
	TimeMage
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

type ClassDefinition struct {
	Name        string
	LayerColors map[string]color.RGBA
	Special1    abilities.Ability
	Special2    abilities.Ability
}
