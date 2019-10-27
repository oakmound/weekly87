package players

import (
	"image/color"
	"strings"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/abilities"
)
// Init to be run after oak setup to get our assets set up
func Init() {
	WarriorsInit()
	mageInit()
	EmptyInit()
	InnKeeperInit()
	classmapping = map[int]*Constructor{
		Empty:     EmptyConstructor,
		InnKeeper: InnKeeperConstructor,
		Spearman:  WarriorConstructors["Spearman"],
		Swordsman: WarriorConstructors["Swordsman"],
		Berserker: WarriorConstructors["Berserker"],
		Paladin:   WarriorConstructors["Paladin"],
		Mage:      mageConstructors["Red"],
		WhiteMage: mageConstructors["White"],
		BlueMage:  mageConstructors["Blue"],
		TimeMage:  mageConstructors["Time"],
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

// Character type enum enum
const (
	Empty = iota
	Swordsman
	Berserker
	Paladin
	Spearman
	Mage
	WhiteMage
	BlueMage
	TimeMage
	InnKeeper
)

var classmapping map[int]*Constructor

// ClassConstructor creates the character types from a PartyMember list
func ClassConstructor(partyComp []PartyMember) []Constructor {
	classes := make([]Constructor, len(partyComp))
	for i, c := range partyComp {
		classes[i] = *classmapping[c.PlayerClass].Copy()
	}
	return classes
}

// ClassDefinition specifies what makes a class special!
type ClassDefinition struct {
	Name        string
	LayerColors map[string]color.RGBA
	Special1    abilities.Ability
	Special2    abilities.Ability
}

// PartyMember information for storage
type PartyMember struct {
	PlayerClass  int
	AccruedValue int
	Name         string
}
