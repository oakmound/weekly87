package section

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/enemies"

	"github.com/oakmound/oak/alg"

	"github.com/oakmound/oak/alg/floatgeom"
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

	// The following section is test code
	testEC := &enemies.Constructor{
		Position:   floatgeom.Point2{400, 400},
		Dimensions: floatgeom.Point2{32, 32},
		Speed:      floatgeom.Point2{-3 * rand.Float64(), -5 * (rand.Float64() - .5)},
		AnimationMap: map[string]render.Modifiable{
			"standRT": render.NewColorBox(32, 32, color.RGBA{255, 125, 0, 255}),
			"standLT": render.NewColorBox(32, 32, color.RGBA{125, 255, 0, 255}),
			"walkRT":  render.NewColorBox(32, 32, color.RGBA{0, 0, 0, 255}),
			"walkLT":  render.NewColorBox(32, 32, color.RGBA{255, 255, 255, 255}),
		},
	}
	e, err := testEC.NewEnemy()
	if err != nil {
		dlog.Error(err)
	} else {
		st.entities = append(st.entities, e)
	}

	if st.sectionsDeep == 1 {
		d := doodads.NewOutDoor(delta < 0)
		d.SetPos(0, 0)
		st.entities = append(st.entities, d)
	} else {
		ch := doodads.NewChest(1)
		ch.SetPos(800, 500)
		st.entities = append(st.entities, ch)
	}

	return st.generate()
}
