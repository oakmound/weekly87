package section

import (
	"math/rand"

	"github.com/oakmound/oak"

	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/enemies"

	"github.com/200sc/go-dist/floatrange"

	"github.com/oakmound/oak/alg"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
)

type Tracker struct {
	start        int64
	sectionsDeep int64
	rng          *rand.Rand
	*compressor
	changes map[int64][]Change
}

func NewTracker(baseSeed int64) *Tracker {
	return &Tracker{
		start:      baseSeed,
		rng:        rand.New(rand.NewSource(baseSeed)),
		compressor: &compressor{},
		changes:    make(map[int64][]Change),
	}
}

// SetDepth is to be used when the chest is picked up, in case there are
// sections that have been generated but they are still off screen.
func (st *Tracker) SetDepth(depth int64) {
	if depth < 1 {
		dlog.Error(`Tracker cannot have depth set to`, depth, `otherwise
			it could fail to create the starting section.`)
		return
	}
	st.sectionsDeep = depth
}

// ShiftDepth is an alternative to SetDepth in case sections aren't being tracked
// outside of the section tracker, but the game knows how many sections ahead
// it has generated from the current one
func (st *Tracker) ShiftDepth(depth int64) {
	st.SetDepth(st.sectionsDeep + depth)
}

func (st *Tracker) AtStart() bool {
	return st.sectionsDeep == 1
}

func (st *Tracker) SectionsDeep() int64 {
	return st.sectionsDeep
}

func (st *Tracker) At() int64 {
	return st.sectionsDeep
}

func (st *Tracker) Prev() *Section {
	return st.Produce(-1)
}

// Next produces another section.
func (st *Tracker) Next() *Section {
	return st.Produce(1)
}

func (st *Tracker) Produce(delta int64) *Section {
	st.sectionsDeep += delta
	st.rng.Seed(st.start + st.sectionsDeep)
	// Section plan:
	// Tiles:
	// A | B
	// - - -
	// C | D
	//
	// Clear Weather
	// A: 3 sections
	// B: 3 sections
	// C: 3 sections
	// D: 3 sections
	// Cloudy weather
	// C: 4 sections
	// B: 4 sections
	// A: 4 sections
	// Stormy weather
	// B: 5 sections
	// C: 5 sections
	// D: 5 sections
	// Snowy weather
	// C: 6 sections
	// B: 6 sections
	// A: 3 Sections
	// Repeat

	// These initial rng calls should make these test sections more distinct
	plan := sectionPlans[((st.sectionsDeep-1)/3)%int64(len(sectionPlans))]
	gWeights := alg.RemainingWeights(plan.groundTileWeights)
	sfWeights := alg.RemainingWeights(plan.surfaceTileWeights)
	skWeights := alg.RemainingWeights(plan.skyTileWeights)

	for x := 0; x < len(st.ground); x++ {
		for y := 0; y < len(st.ground[x]); y++ {
			choice := alg.WeightedChooseOne(gWeights)
			t := plan.groundTiles[choice]
			st.ground[x][y] = t.Copy()
		}
	}
	for x := 0; x < len(st.wall); x++ {
		for y := 0; y < len(st.wall[x]); y++ {
			var t render.Modifiable
			if y == len(st.wall[x])-1 {
				choice := alg.WeightedChooseOne(sfWeights)
				t = plan.surfaceTiles[choice]

			} else {
				choice := alg.WeightedChooseOne(skWeights)
				t = plan.skyTiles[choice]
			}
			st.wall[x][y] = t.Copy()
		}
	}

	fieldX := floatrange.NewLinear(0, float64(oak.ScreenWidth))
	fieldY := floatrange.NewLinear(float64(oak.ScreenHeight)*1/3, float64(oak.ScreenHeight)-64)

	enemyDist := alg.RemainingWeights(plan.enemyDistribution[:])

	if !(st.sectionsDeep == 1 && delta > 0) {
		for i := 0; i < plan.enemyCount.Poll(); i++ {
			typ := alg.WeightedChooseOne(enemyDist)
			cs := enemies.Constructors[typ]
			e, err := cs.NewEnemy(st.sectionsDeep, int64(len(st.entities)))
			if delta < 0 {
				e.RunBackwards()
			}
			dlog.ErrorCheck(err)
			e.SetPos(fieldX.Poll(), fieldY.Poll())
			st.entities = append(st.entities, e)
		}
	}

	if st.sectionsDeep == 1 {
		d := doodads.NewOutDoor(delta < 0)
		d.SetPos(0, float64(oak.ScreenHeight-10)*1/3)
		st.entities = append(st.entities, d)
	} else if st.sectionsDeep > 2 {
		for i := 0; i < plan.chestCount.Poll(); i++ {
			ch := doodads.NewChest(int64(plan.chestRange.Poll()))
			ch.SetPos(fieldX.Poll(), fieldY.Poll())
			st.entities = append(st.entities, ch)
		}
	}

	// for i, e := range st.entities {
	// 	e.SetIdx(i)
	// }

	newSection := st.generate()
	for _, c := range st.changes[newSection.id] {
		newSection.ApplyChange(c)
	}

	return newSection
}

func (st *Tracker) UpdateHistory(sectionID int64, change Change) {
	st.changes[sectionID] = append(st.changes[sectionID], change)
}
